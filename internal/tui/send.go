package tui

import (
	"os"
	"fmt"
	"net"
	"math"
	"time"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/abdfnx/tran/tools"
	"github.com/abdfnx/tran/models"
	"github.com/abdfnx/tran/constants"
	"github.com/abdfnx/tran/core/sender"
	tea "github.com/charmbracelet/bubbletea"
)

func HandleSendCommand(programOptions models.TranOptions, fileNames []string) {
	// communicate ui updates on this channel between senderClient and HandleSendCommand
	uiCh := make(chan sender.UIUpdate)
	// initialize a senderClient with a UI
	senderClient := sender.WithUI(sender.NewSender(programOptions), uiCh)
	// initialize and start sender-UI
	senderUI := NewSenderUI()
	// clean up temporary files previously created by this command
	tools.RemoveTemporaryFiles(constants.SEND_TEMP_FILE_NAME_PREFIX)

	go initSenderUI(senderUI)
	time.Sleep(constants.START_PERIOD)
	go listenForSenderUIUpdates(senderUI, uiCh)

	closeFileCh := make(chan *os.File)
	senderReadyCh := make(chan bool, 1)
	// read, archive and compress files in parallel
	go prepareFiles(senderClient, senderUI, fileNames, senderReadyCh, closeFileCh)

	// initiate communications with tranx-server
	startServerCh := make(chan sender.ServerOptions)
	relayCh := make(chan *websocket.Conn)
	passCh := make(chan models.Password)
	go initiateSenderTranxCommunication(senderClient, senderUI, passCh, startServerCh, senderReadyCh, relayCh)
	senderUI.Send(PasswordMsg{Password: string(<-passCh)})

	// keeps program alive until finished
	doneCh := make(chan bool)
	// attach server to senderClient
	senderClient = sender.WithServer(senderClient, <-startServerCh)

	go startDirectCommunicationServer(senderClient, senderUI, doneCh)
	// prepare a fallback to relay communications through tranx if direct communications unavailble
	prepareRelayCommunicationFallback(senderClient, senderUI, relayCh, doneCh)

	<-doneCh
	senderUI.Send(FinishedMsg{})
	tempFile := <-closeFileCh
	os.Remove(tempFile.Name())
	tempFile.Close()
	GracefulUIQuit()
}

func initSenderUI(senderUI *tea.Program) {
	if err := senderUI.Start(); err != nil {
		fmt.Println("Error initializing UI", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func listenForSenderUIUpdates(senderUI *tea.Program, uiCh chan sender.UIUpdate) {
	latestProgress := 0
	for uiUpdate := range uiCh {
		// make sure progress is 100 if connection is to be closed
		if uiUpdate.State == sender.WaitForCloseMessage {
			latestProgress = 100
			senderUI.Send(ProgressMsg{Progress: 1})
			continue
		}

		// limit progress update ui-send events
		newProgress := int(math.Ceil(100 * float64(uiUpdate.Progress)))
		if newProgress > latestProgress {
			latestProgress = newProgress
			senderUI.Send(ProgressMsg{Progress: uiUpdate.Progress})
		}
	}
}

func prepareFiles(senderClient *sender.Sender, senderUI *tea.Program, fileNames []string, readyCh chan bool, closeFileCh chan *os.File) {
	files, err := tools.ReadFiles(fileNames)

	if err != nil {
		senderUI.Send(ErrorMsg{Message: "Error reading files."})
		GracefulUIQuit()
	}

	uncompressedFileSize, err := tools.FilesTotalSize(files)
	if err != nil {
		senderUI.Send(ErrorMsg{Message: "Error during file preparation."})
		GracefulUIQuit()
	}

	senderUI.Send(FileInfoMsg{FileNames: fileNames, Bytes: uncompressedFileSize})

	tempFile, fileSize, err := tools.ArchiveAndCompressFiles(files)
	for _, file := range files {
		file.Close()
	}

	if err != nil {
		senderUI.Send(ErrorMsg{Message: "Error compressing files."})
		GracefulUIQuit()
	}

	sender.WithPayload(senderClient, tempFile, fileSize)
	senderUI.Send(FileInfoMsg{FileNames: fileNames, Bytes: fileSize})
	readyCh <- true
	senderUI.Send(ReadyMsg{})
	closeFileCh <- tempFile
}

func initiateSenderTranxCommunication(senderClient *sender.Sender, senderUI *tea.Program, passCh chan models.Password,
	startServerCh chan sender.ServerOptions, readyCh chan bool, relayCh chan *websocket.Conn) {
	err := senderClient.ConnectToTranx(
		senderClient.TranxAddress(), senderClient.TranxPort(), passCh, startServerCh, readyCh, relayCh)

	if err != nil {
		senderUI.Send(ErrorMsg{Message: "Failed to communicate with tranx server."})
		GracefulUIQuit()
	}
}

func startDirectCommunicationServer(senderClient *sender.Sender, senderUI *tea.Program, doneCh chan bool) {
	if err := senderClient.StartServer(); err != nil {
		senderUI.Send(ErrorMsg{Message: fmt.Sprintf("Something went wrong during file transfer: %e", err)})
		GracefulUIQuit()
	}

	doneCh <- true
}

func prepareRelayCommunicationFallback(senderClient *sender.Sender, senderUI *tea.Program, relayCh chan *websocket.Conn, doneCh chan bool) {
	if relayWsConn, closed := <-relayCh; closed {
		// start transferring to the tranx-relay
		go func() {
			if err := senderClient.Transfer(relayWsConn); err != nil {
				senderUI.Send(ErrorMsg{Message: fmt.Sprintf("Something went wrong during file transfer: %e", err)})
				GracefulUIQuit()
			}

			doneCh <- true
		}()
	}
}

func ValidateTranxAddress() error {
	address := net.ParseIP(constants.DEFAULT_ADDRESS)
	err := tools.ValidateHostname(constants.DEFAULT_ADDRESS)

	// neither a valid IP nor a valid hostname was provided
	if (address == nil) && err != nil {
		return errors.New("invalid IP or hostname provided")
	}

	return nil
}

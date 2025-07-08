package tui

import (
	"os"
	"fmt"
	"math"
	"time"

	"github.com/gorilla/websocket"
	"github.com/abdfnx/tran/tools"
	"github.com/abdfnx/tran/models"
	"github.com/abdfnx/tran/constants"
	"github.com/abdfnx/tran/core/receiver"
	"github.com/abdfnx/tran/models/protocol"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleReceiveCommand is the receive application.
func HandleReceiveCommand(programOptions models.TranOptions, password string) {
	// communicate ui updates on this channel between receiverClient and handleReceiveCommand
	uiCh := make(chan receiver.UIUpdate)
	// initialize a receiverClient with a UI
	receiverClient := receiver.WithUI(receiver.NewReceiver(programOptions), uiCh)
	// initialize and start receiver-UI
	receiverUI := NewReceiverUI()
	// clean up temporary files previously created by this command
	tools.RemoveTemporaryFiles(constants.RECEIVE_TEMP_FILE_NAME_PREFIX)

	go initReceiverUI(receiverUI)
	time.Sleep(constants.START_PERIOD)
	go listenForReceiverUIUpdates(receiverUI, uiCh)

	parsedPassword, err := tools.ParsePassword(password)
	if err != nil {
		receiverUI.Send(ErrorMsg{Message: "Error parsing password, make sure you entered a correctly formatted password (e.g. 1-gamma-ray-quasar)."})
		GracefulUIQuit()
	}

	// initiate communications with tranx-server
	wsConnCh := make(chan *websocket.Conn)
	go initiateReceiverTranxCommunication(receiverClient, receiverUI, parsedPassword, wsConnCh)

	// keeps program alive until finished
	doneCh := make(chan bool)
	// start receiving files
	go startReceiving(receiverClient, receiverUI, <-wsConnCh, doneCh)

	// wait for shut down to render final UI
	<-doneCh
	GracefulUIQuit()
}

func initReceiverUI(receiverUI *tea.Program) {
	go func() {
		if err := receiverUI.Start(); err != nil {
			fmt.Println("Error initializing UI", err)
			os.Exit(1)
		}

		os.Exit(0)
	}()
}

func listenForReceiverUIUpdates(receiverUI *tea.Program, uiCh chan receiver.UIUpdate) {
	latestProgress := 0

	for uiUpdate := range uiCh {
		// limit progress update ui-send events
		newProgress := int(math.Ceil(100 * float64(uiUpdate.Progress)))
		if newProgress > latestProgress {
			latestProgress = newProgress
			receiverUI.Send(ProgressMsg{Progress: uiUpdate.Progress})
		}
	}
}

func initiateReceiverTranxCommunication(receiverClient *receiver.Receiver, receiverUI *tea.Program, password models.Password, connectionCh chan *websocket.Conn) {
	wsConn, err := receiverClient.ConnectToTranx(receiverClient.TranxAddress(), receiverClient.TranxPort(), password)
	if err != nil {
		receiverUI.Send(ErrorMsg{Message: "Something went wrong during connection-negotiation (did you enter the correct password?)"})
		GracefulUIQuit()
	}

	receiverUI.Send(FileInfoMsg{Bytes: receiverClient.PayloadSize()})
	connectionCh <- wsConn
}

func startReceiving(receiverClient *receiver.Receiver, receiverUI *tea.Program, wsConnection *websocket.Conn, doneCh chan bool) {
	tempFile, err := os.CreateTemp(os.TempDir(), constants.RECEIVE_TEMP_FILE_NAME_PREFIX)

	if err != nil {
		receiverUI.Send(ErrorMsg{Message: "Something went wrong when creating the received file container."})
		GracefulUIQuit()
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// start receiving files from sender
	err = receiverClient.Receive(wsConnection, tempFile)
	if err != nil {
		receiverUI.Send(ErrorMsg{Message: "Something went wrong during file transfer."})
		GracefulUIQuit()
	}

	if receiverClient.UsedRelay() {
		wsConnection.WriteJSON(protocol.TranxMessage{Type: protocol.ReceiverToTranxClose})
	}

	// reset file position for reading
	tempFile.Seek(0, 0)

	// read received bytes from tmpFile
	receivedFileNames, decompressedSize, err := tools.DecompressAndUnarchiveBytes(tempFile)
	if err != nil {
		receiverUI.Send(ErrorMsg{Message: "Something went wrong when expanding the received files."})
		GracefulUIQuit()
	}

	receiverUI.Send(FinishedMsg{Files: receivedFileNames, PayloadSize: decompressedSize})
	doneCh <- true
}

package tui

import (
	"fmt"
	"sort"
	"time"
	"strings"

	"github.com/abdfnx/tran/constants"
	"github.com/charmbracelet/bubbles/spinner"
)

type UIUpdate struct {
	Progress float32
}

type FileInfoMsg struct {
	FileNames []string
	Bytes     int64
}

type ErrorMsg struct {
	Message string
}

type ProgressMsg struct {
	Progress float32
}

type FinishedMsg struct {
	Files       []string
	PayloadSize int64
}

var WaitingSpinner = spinner.Dot

var CompressingSpinner = spinner.Globe

var TransferSpinner = spinner.Spinner{
	Frames: []string{"»  ", "»» ", "»»»", "   "},
	FPS:    time.Millisecond * 400,
}

var ReceivingSpinner = spinner.Spinner{
	Frames: []string{"   ", "  «", " ««", "«««"},
	FPS:    time.Second / 2,
}

func TopLevelFilesText(fileNames []string) string {
	// parse top level file names and attach number of subfiles in them
	topLevelFileChildren := make(map[string]int)

	for _, f := range fileNames {
		fileTopPath := strings.Split(f, "/")[0]

		subfileCount, wasPresent := topLevelFileChildren[fileTopPath]

		if wasPresent {
			topLevelFileChildren[fileTopPath] = subfileCount + 1
		} else {
			topLevelFileChildren[fileTopPath] = 0
		}
	}

	// read map into formatted strings
	var topLevelFilesText []string

	for fileName, subFileCount := range topLevelFileChildren {
		formattedFileName := fileName

		if subFileCount > 0 {
			formattedFileName = fmt.Sprintf("%s (%d subfiles)", fileName, subFileCount)
		}

		topLevelFilesText = append(topLevelFilesText, formattedFileName)
	}

	sort.Strings(topLevelFilesText)

	return strings.Join(topLevelFilesText, ", ")
}

func GracefulUIQuit() {
	time.Sleep(constants.SHUTDOWN_PERIOD)
}

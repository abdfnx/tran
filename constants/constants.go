package constants

import (
	"fmt"
	"time"
	"runtime"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/progress"
)

const DEFAULT_ADDRESS = "167.71.65.96"
const DEFAULT_PORT = 80

const MAX_CHUNK_BYTES = 1e6
const MAX_SEND_CHUNKS = 2e8

const RECEIVER_CONNECT_TIMEOUT time.Duration = 5 * time.Minute

const SEND_TEMP_FILE_NAME_PREFIX = "tran-send-tmp"
const RECEIVE_TEMP_FILE_NAME_PREFIX = "tran-receive-tmp"

const (
	PrimaryBoxActive = iota
	SecondaryBoxActive
	ThirdBoxActive
)

const (
	StatusBarHeight      = 1
	BoxPadding           = 1
	EllipsisStyle        = "..."
	FileSizeLoadingStyle = "---"
)

var BoldTextStyle = lipgloss.NewStyle().Bold(true)

var Colors = map[string]lipgloss.Color{
	"black": "#000000",
}

const (
	PADDING                 = 2
	MAX_WIDTH               = 80
	PRIMARY_COLOR           = "#1E90FF"
	SECONDARY_COLOR         = "#1E6AFF"
	DARK_GRAY_COLOR         = "#3c3836"
	START_PERIOD            = 1 * time.Millisecond
	SHUTDOWN_PERIOD         = 1000 * time.Millisecond
)

var QuitKeys = []string{"q", "esc"}
var PadText = strings.Repeat(" ", PADDING)
var QuitCommandsHelpText = HelpStyle(fmt.Sprintf("(press one of [%s] keys to quit from tran)", (strings.Join(QuitKeys, ", "))))
var ProgressBar = progress.NewModel(progress.WithGradient(SECONDARY_COLOR, PRIMARY_COLOR))

var baseStyle = lipgloss.NewStyle()
var InfoStyle = baseStyle.Copy().Foreground(lipgloss.Color(PRIMARY_COLOR)).Render
var HelpStyle = baseStyle.Copy().Foreground(lipgloss.Color(DARK_GRAY_COLOR)).Render
var ItalicText = baseStyle.Copy().Italic(true).Render
var BoldText = baseStyle.Copy().Bold(true).Render

func CtrlKey() string {
	// if os is macos, then return "⌘"
	if runtime.GOOS == "darwin" {
		return "⌘"
	} else {
		return "ctrl"
	}
}

func AltKey() string {
	// if os is macos, then return "⌥"
	if runtime.GOOS == "darwin" {
		return "⌥"
	} else {
		return "alt"
	}
}

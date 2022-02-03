package constants

import (
	"time"
	"runtime"

	"github.com/charmbracelet/lipgloss"
)

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
	START_PERIOD            = 1 * time.Millisecond
	SHUTDOWN_PERIOD         = 1000 * time.Millisecond
)

const DEFAULT_ADDRESS = "167.71.65.96"
const DEFAULT_PORT = 80

func CtrlKey() string {
	// if os is macos, then return "⌘"
	if runtime.GOOS == "darwin" {
		return "⌘"
	} else {
		return "ctrl"
	}
}

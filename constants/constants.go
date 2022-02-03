package constants

import (
	"runtime"
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

package checker

import (
	"fmt"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/abdfnx/looker"
	"github.com/abdfnx/tran/api"
	"github.com/abdfnx/tran/cmd/factory"
	"github.com/abdfnx/tran/internal/config"
)

func Check(buildVersion string) {
	cmdFactory := factory.New()
	stderr := cmdFactory.IOStreams.ErrOut
	cfg := config.GetConfig()

	latestVersion := api.GetLatest()
	isFromHomebrewTap := isUnderHomebrew()
	isFromGo := isUnderGo()
	isFromUsrBinDir := isUnderUsr()
	isFromGHCLI := isUnderGHCLI()
	isFromAppData := isUnderAppData()

	var command = func() string {
		if isFromHomebrewTap {
			return "brew upgrade tran"
		} else if isFromGo {
			return "go get -u github.com/abdfnx/tran"
		} else if isFromUsrBinDir {
			return "curl -fsSL https://cutt.ly/tran-cli | bash"
		} else if isFromGHCLI {
			return "gh extention upgrade tran"
		} else if isFromAppData {
			return "iwr -useb https://cutt.ly/tran-win | iex"
		}

		return ""
	}

	if buildVersion != latestVersion && cfg.Tran.ShowUpdates != false {
		fmt.Fprintf(stderr, "%s %s → %s\n",
		ansi.Color("There's a new version of ", "yellow") + ansi.Color("tran", "cyan") + ansi.Color(" is avalaible:", "yellow"),
		ansi.Color(buildVersion, "cyan"),
		ansi.Color(latestVersion, "cyan"))

		if command() != "" {
			fmt.Fprintf(stderr, ansi.Color("To upgrade, run: %s\n", "yellow"), ansi.Color(command(), "black:white"))
		}
	}
}

var tranExe, _ = looker.LookPath("tran")

func isUnderHomebrew() bool {
	return strings.Contains(tranExe, "brew")
}

func isUnderGo() bool {
	return strings.Contains(tranExe, "go")
}

func isUnderUsr() bool {
	return strings.Contains(tranExe, "usr")
}

func isUnderAppData() bool {
	return strings.Contains(tranExe, "AppData")
}

func isUnderGHCLI() bool {
	return strings.Contains(tranExe, "gh")
}

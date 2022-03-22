package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	surveyCore "github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/mgutz/ansi"
	"github.com/scmn-dev/tran/app/checker"
	"github.com/scmn-dev/tran/cmd/factory"
	"github.com/scmn-dev/tran/cmd/tran"
	"github.com/scmn-dev/tran/tools"
	"github.com/spf13/cobra"
)

var (
	version string
	buildDate string
)

type exitCode int

const (
	exitOK     exitCode = 0
	exitError  exitCode = 1
	exitCancel exitCode = 2
)

func main() {
	code := mainRun()
	os.Exit(int(code))
}

func mainRun() exitCode {
	runtime.LockOSThread()

	cmdFactory := factory.New()
	hasDebug := os.Getenv("DEBUG") != ""
	stderr := cmdFactory.IOStreams.ErrOut

	if !cmdFactory.IOStreams.ColorEnabled() {
		surveyCore.DisableColor = true
	} else {
		surveyCore.TemplateFuncsWithColor["color"] = func(style string) string {
			switch style {
				case "white":
					if cmdFactory.IOStreams.ColorSupport256() {
						return fmt.Sprintf("\x1b[%d;5;%dm", 38, 242)
					}

					return ansi.ColorCode("default")

				default:
					return ansi.ColorCode(style)
			}
		}
	}

	if len(os.Args) > 1 && os.Args[1] != "" {
		cobra.MousetrapHelpText = ""
	}

	RootCmd := tran.Execute(cmdFactory, version, buildDate)

	if cmd, err := RootCmd.ExecuteC(); err != nil {
		if err == tools.SilentError {
			return exitError
		} else if tools.IsUserCancellation(err) {
			if errors.Is(err, terminal.InterruptErr) {
				fmt.Fprint(stderr, "\n")
			}

			return exitCancel
		}

		tools.PrintError(stderr, err, cmd, hasDebug)

		return exitError
	}

	if tran.HasFailed() {
		return exitError
	}

	if len(os.Args) > 1 && os.Args[1] != "tran" {
		checker.Check(version)
	}

	return exitOK
}

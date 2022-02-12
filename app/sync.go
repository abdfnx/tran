package app

import (
	"fmt"
	"log"
	"time"
	"runtime"

	"github.com/abdfnx/gosh"
	"github.com/spf13/cobra"
	"github.com/abdfnx/gh/utils"
	"github.com/briandowns/spinner"
	"github.com/MakeNowJust/heredoc"
	"github.com/abdfnx/tran/constants"
	git_config "github.com/david-tomson/tran-git"
)

var username = git_config.GitConfig()

var (
	NewCmdStart = &cobra.Command{
		Use:   "start",
		Aliases: []string{"."},
		Example: "tran sync start",
		Short: "Start sync your tran config.",
		Run: func(cmd *cobra.Command, args []string) {
			if username != ":username" {
				exCmd := "echo '# My tran config - " + username + "\n\n## Clone\n\n```\ntran sync clone\n```\n\n**for more about sync command, run `tran sync -h`**' >> $HOME/.config/tran/README.md"

				gosh.Run(exCmd)
				gosh.RunMulti(constants.Start_ml(), constants.Start_w())
			} else {
				utils.AuthMessage()
			}
		},
	}

	NewCmdClone = &cobra.Command{
		Use:   "clone",
		Aliases: []string{"cn"},
		Short: CloneHelp(),
		Run: func(cmd *cobra.Command, args []string) {
			if username != ":username" {
				gosh.RunMulti(constants.Clone_ml(), constants.Clone_w())
				gosh.RunMulti(constants.Clone_check_ml(), constants.Clone_check_w())
			} else {
				utils.AuthMessage()
			}
		},
	}

	NewCmdPush = &cobra.Command{
		Use:   "push",
		Aliases: []string{"ph"},
		Short: "Push the new changes in tran config file.",
		Run: func(cmd *cobra.Command, args []string) {
			if username != ":username" {
				gosh.RunMulti(constants.Push_ml(), constants.Push_w())
			} else {
				utils.AuthMessage()
			}
		},
	}

	NewCmdPull = &cobra.Command{
		Use:   "pull",
		Aliases: []string{"pl"},
		Short: PullHelp(),
		Run: func(cmd *cobra.Command, args []string) {
			if username != ":username" {
				gosh.RunMulti(constants.Pull_ml(), constants.Pull_w())
			} else {
				utils.AuthMessage()
			}
		},
	}

	FetchX = &cobra.Command{
		Use:   "fetchx",
		Short: "Special command for windows",
		Run: func(cmd *cobra.Command, args []string) {
			if username != ":username" {
				if runtime.GOOS == "windows" {
					gosh.PowershellCommand(constants.Clone_w())
				} else {
					fmt.Println("This command isn't avaliable for this platform")
				}
			} else {
				utils.AuthMessage()
			}
		},
	}
)

func Sync() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sync <command>",
		Short: "Sync your tran config file.",
		Long:  SyncHelp(),
		Example: heredoc.Doc(`
			tran sync start
			tran sync clone
		`),
	}

	cmd.AddCommand(
		NewCmdStart,
		NewCmdClone,
		NewCmdPush,
		NewCmdPull,
		FetchX,
	)

	return cmd
}

const tranConfigPath string = "/.tran"

func PullHelp() string {
	return git_config.GitConfigWithMsg("Pull the new changes from ", tranConfigPath)
}

func SyncHelp() string {
	return git_config.GitConfigWithMsg("Sync your config file, by create a private repo at ", tranConfigPath)
}

func CloneHelp() string {
	return git_config.GitConfigWithMsg("Clone your .tran from your private repo at https://github.com/", tranConfigPath)
}

func PushSync() {
	const Syncing string = " 📮 Syncing..."

	if runtime.GOOS == "windows" {
		err, out, errout := gosh.PowershellOutput(
		`
			$directoyPath = "~/.config/tran/.git"

			if (Test-Path -path $directoyPath) {
				Write-Host "Reading from .tran folder..."
			}
		`)

		fmt.Print(out)

		if err != nil {
			log.Printf("error: %v\n", err)
			fmt.Print(errout)
		} else if out != "" {
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = Syncing
			s.Start()

			gosh.PowershellCommand(constants.Push_w())

			s.Stop()
		}
	} else {
		err, out, errout := gosh.ShellOutput(
		`
			if [ -d ~/.config/tran/.git ]; then
				echo "📖 Reading from .tran folder..."
			fi
		`)

		fmt.Print(out)

		if err != nil {
			log.Printf("error: %v\n", err)
			fmt.Print(errout)
		} else if out != "" {
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = Syncing
			s.Start()

			gosh.ShellCommand(constants.Push_ml())

			s.Stop()
		}
	}
}


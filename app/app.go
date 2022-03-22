package app

import (
	"log"

	"github.com/abdfnx/gh/pkg/cmd/factory"
	"github.com/scmn-dev/tran/constants"
	"github.com/scmn-dev/tran/models"
	"github.com/scmn-dev/tran/tools"
	"github.com/scmn-dev/tran/tui"
	"github.com/spf13/cobra"
)

var NewSendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send files/directories to remote",
	Long:  "Send files/directories to remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		tools.RandomSeed()

		err := tui.ValidateTranxAddress()

		if err != nil {
			log.Fatal(err)
		}

		tui.HandleSendCommand(models.TranOptions{
			TranxAddress: constants.DEFAULT_ADDRESS,
			TranxPort:    constants.DEFAULT_PORT,
		}, args)

		return nil
	},
}

var NewReceiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive files/directories from remote",
	Long:  "Receive files/directories from remote",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := tui.ValidateTranxAddress()

		if err != nil {
			return err
		}

		tui.HandleReceiveCommand(models.TranOptions{
			TranxAddress: constants.DEFAULT_ADDRESS,
			TranxPort:    constants.DEFAULT_PORT,
		}, args[0])

		return nil
	},
}

var NewAuthCmd = Auth(factory.New())
var NewGHConfigCmd = GHConfig(factory.New())
var NewGHRepoCmd = Repo(factory.New())

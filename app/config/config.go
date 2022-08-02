package config

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"path/filepath"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/abdfnx/tran/dfs"
)

var (
	homeDir, _ = dfs.GetHomeDirectory()
	tranConfigPath = filepath.Join(homeDir, ".tran", "tran.yml")
	tranConfig, err = ioutil.ReadFile(tranConfigPath)
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure tran",
		Long:  "Configure tran, including setting up tran editor, etc.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cmd.AddCommand(NewConfigSetCmd)
	cmd.AddCommand(NewConfigGetCmd)
	cmd.AddCommand(NewConfigListCmd)

	return cmd
}

var NewConfigSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Update tran configuration",
	Long:  "Update tran configuration, such as editor, show updates, etc.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err != nil {
			return err
		}

		viper.SetConfigType("yaml")

		viper.ReadConfig(bytes.NewBuffer(tranConfig))

		// set new key value but keep existing values
		viper.Set("config." + args[0], args[1])

		// write config to file
		err := viper.WriteConfigAs(tranConfigPath)

		if err != nil {
			return err
		}

		fmt.Println(ansi.Color("Updated tran configuration", "green"))

		return nil
	},
}

var NewConfigGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get tran configuration",
	Long:  "Get tran configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err != nil {
			return err
		}

		viper.SetConfigType("yaml")

		viper.ReadConfig(bytes.NewBuffer(tranConfig))
		
		fmt.Println(viper.Get("config." + args[0]))

		return nil
	},
}

var NewConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tran configuration",
	Long:  "List tran configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err != nil {
			return err
		}

		viper.SetConfigType("yaml")

		viper.ReadConfig(bytes.NewBuffer(tranConfig))
		
		// get the config
		config := viper.GetStringMap("config")

		for k, v := range config {
			fmt.Println(k + " =", v)
		}

		return nil
	},
}

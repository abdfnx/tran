package app

import (
	"github.com/spf13/cobra"
	"github.com/abdfnx/gh/pkg/cmdutil"
	aCmd "github.com/abdfnx/gh/pkg/cmd/auth"
	cCmd "github.com/abdfnx/gh/pkg/cmd/gh-config"
)

func Auth(f *cmdutil.Factory) *cobra.Command {
	cmd := aCmd.NewCmdAuth(f)
	return cmd
}

func GHConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := cCmd.NewCmdConfig(f)
	return cmd
}

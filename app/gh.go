package app

import (
	"github.com/spf13/cobra"
	"github.com/abdfnx/gh/pkg/cmdutil"
	aCmd "github.com/abdfnx/gh/pkg/cmd/auth"
	cCmd "github.com/abdfnx/gh/pkg/cmd/gh-config"
	rCmd "github.com/abdfnx/gh/pkg/cmd/gh-repo"
	"github.com/abdfnx/gh/context"
	"github.com/abdfnx/gh/api"
	"github.com/abdfnx/gh/core/ghrepo"
)

func Auth(f *cmdutil.Factory) *cobra.Command {
	cmd := aCmd.NewCmdAuth(f)
	return cmd
}

func GHConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := cCmd.NewCmdConfig(f)
	return cmd
}

func Repo(f *cmdutil.Factory) *cobra.Command {
	repoResolvingCmdFactory := *f
	repoResolvingCmdFactory.BaseRepo = resolvedBaseRepo(f)

	cmd := rCmd.NewCmdRepo(&repoResolvingCmdFactory)

	return cmd
}

func resolvedBaseRepo(f *cmdutil.Factory) func() (ghrepo.Interface, error) {
	return func() (ghrepo.Interface, error) {
		httpClient, err := f.HttpClient()
		if err != nil {
			return nil, err
		}

		apiClient := api.NewClientFromHTTP(httpClient)

		remotes, err := f.Remotes()
		if err != nil {
			return nil, err
		}

		repoContext, err := context.ResolveRemotesToRepos(remotes, apiClient, "")
		if err != nil {
			return nil, err
		}

		baseRepo, err := repoContext.BaseRepo(f.IOStreams)
		if err != nil {
			return nil, err
		}

		return baseRepo, nil
	}
}

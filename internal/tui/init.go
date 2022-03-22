package tui

import (
	"os"
	"log"
	"strings"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/abdfnx/tran/dfs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

// Init initializes the UI and sets up initial data.
func (b Bubble) Init() tea.Cmd {
	var cmds []tea.Cmd
	startDir := viper.GetString("start-dir")

	switch {
		case startDir != "":
			_, err := os.Stat(startDir)
			if err != nil {
				return nil
			}

			if strings.HasPrefix(startDir, dfs.RootDirectory) {
				cmds = append(cmds, b.updateDirectoryListingCmd(startDir))
			} else {
				path, err := os.Getwd()

				if err != nil {
					log.Fatal(err)
				}

				filePath := filepath.Join(path, startDir)

				cmds = append(cmds, b.updateDirectoryListingCmd(filePath))
			}

		case b.appConfig.Tran.StartDir == dfs.HomeDirectory:
			homeDir, err := dfs.GetHomeDirectory()
			if err != nil {
				log.Fatal(err)
			}

			cmds = append(cmds, b.updateDirectoryListingCmd(homeDir))

		default:
			cmds = append(cmds, b.updateDirectoryListingCmd(b.appConfig.Tran.StartDir))
	}

	cmds = append(cmds, spinner.Tick)
	cmds = append(cmds, textinput.Blink)

	return tea.Batch(cmds...)
}

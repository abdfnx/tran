package tui

import (
	"os"
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/muesli/termenv"
	"github.com/abdfnx/tran/dfs"
	"github.com/abdfnx/tran/constants"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/abdfnx/tran/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
)

// checkPrimaryViewportBounds handles wrapping of the filetree and scrolling of the viewport.
func (b *Bubble) checkPrimaryViewportBounds() {
	top := b.primaryViewport.YOffset
	bottom := b.primaryViewport.Height + b.primaryViewport.YOffset - 1

	if b.treeCursor < top {
		b.primaryViewport.LineUp(1)
	} else if b.treeCursor > bottom {
		b.primaryViewport.LineDown(1)
	}

	if b.treeCursor > len(b.treeFiles)-1 {
		b.primaryViewport.GotoTop()
		b.treeCursor = 0
	} else if b.treeCursor < top {
		b.primaryViewport.GotoBottom()
		b.treeCursor = len(b.treeFiles) - 1
	}
}

// handleKeys handles all keypresses.
func (b *Bubble) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	cfg := config.GetConfig()

	// Jump to top of box.
	if msg.String() == "T" {
		if !b.showCommandInput && b.activeBox == constants.PrimaryBoxActive && !b.showBoxSpinner {
			b.treeCursor = 0
			b.primaryViewport.GotoTop()
			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
		}

		if !b.showCommandInput && b.activeBox == constants.SecondaryBoxActive {
			b.secondaryViewport.GotoTop()
		}

		if !b.showCommandInput && b.activeBox == constants.ThirdBoxActive {
			b.thirdViewport.GotoTop()
		}

		return nil
	}

	switch {
		case key.Matches(msg, b.keyMap.Quit):
			return tea.Quit

		case key.Matches(msg, b.keyMap.Down):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor++
				b.checkPrimaryViewportBounds()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

		case key.Matches(msg, b.keyMap.Up):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor--
				b.checkPrimaryViewportBounds()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

		case key.Matches(msg, b.keyMap.Left):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.showFilesOnly = false
				b.showDirectoriesOnly = false
				b.foundFilesPaths = nil
				workingDirectory, err := dfs.GetWorkingDirectory()

				if err != nil {
					return b.handleErrorCmd(err)
				}

				return b.updateDirectoryListingCmd(
					filepath.Join(workingDirectory, dfs.PreviousDirectory),
				)
			}

		case key.Matches(msg, b.keyMap.Right):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile, err := b.treeFiles[b.treeCursor].Info()

				if err != nil {
					return b.handleErrorCmd(err)
				}

				switch {
					case selectedFile.IsDir():
						currentDir, err := dfs.GetWorkingDirectory()
						if err != nil {
							return b.handleErrorCmd(err)
						}

						directoryToOpen := filepath.Join(currentDir, selectedFile.Name())

						if len(b.foundFilesPaths) > 0 {
							directoryToOpen = b.foundFilesPaths[b.treeCursor]
						}

						return b.updateDirectoryListingCmd(directoryToOpen)

					case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
						symlinkFile, err := os.Readlink(selectedFile.Name())

						if err != nil {
							return b.handleErrorCmd(err)
						}

						fileInfo, err := os.Stat(symlinkFile)

						if err != nil {
							return b.handleErrorCmd(err)
						}

						if fileInfo.IsDir() {
							currentDir, err := dfs.GetWorkingDirectory()
							if err != nil {
								return b.handleErrorCmd(err)
							}

							return b.updateDirectoryListingCmd(filepath.Join(currentDir, fileInfo.Name()))
						}

						return b.readFileContentCmd(
							fileInfo.Name(),
							b.secondaryViewport.Width,
						)

					default:
						fileToRead := selectedFile.Name()

						if len(b.foundFilesPaths) > 0 {
							fileToRead = b.foundFilesPaths[b.treeCursor]
						}

						return b.readFileContentCmd(
							fileToRead,
							b.secondaryViewport.Width,
						)
				}
			}

		case key.Matches(msg, b.keyMap.View):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				selectedFile, err := b.treeFiles[b.treeCursor].Info()

				if err != nil {
					return b.handleErrorCmd(err)
				}

				switch {
					case selectedFile.IsDir():
						return b.viewDirectoryListingCmd(selectedFile.Name())

					case selectedFile.Mode()&os.ModeSymlink == os.ModeSymlink:
						symlinkFile, err := os.Readlink(selectedFile.Name())

						if err != nil {
							return b.handleErrorCmd(err)
						}

						fileInfo, err := os.Stat(symlinkFile)
						if err != nil {
							return b.handleErrorCmd(err)
						}

						if fileInfo.IsDir() {
							return b.viewDirectoryListingCmd(fileInfo.Name())
						}

					default:
						return nil
				}
			}

		case key.Matches(msg, b.keyMap.GotoBottom):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = len(b.treeFiles) - 1
				b.primaryViewport.GotoBottom()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == constants.SecondaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.secondaryViewport.GotoBottom()
			}

		case key.Matches(msg, b.keyMap.HomeShortcut):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.fileSizes = nil
				homeDir, err := dfs.GetHomeDirectory()

				if err != nil {
					return b.handleErrorCmd(err)
				}

				return b.updateDirectoryListingCmd(homeDir)
			}

		case key.Matches(msg, b.keyMap.RootShortcut):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.treeCursor = 0
				b.fileSizes = nil

				return b.updateDirectoryListingCmd(dfs.RootDirectory)
			}

		case key.Matches(msg, b.keyMap.ToggleHidden):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.showHiddenFiles = !b.showHiddenFiles

				switch {
					case b.showDirectoriesOnly:
						return b.getDirectoryListingByTypeCmd(dfs.DirectoriesListingType)

					case b.showFilesOnly:
						return b.getDirectoryListingByTypeCmd(dfs.FilesListingType)

					default:
						return b.updateDirectoryListingCmd(dfs.CurrentDirectory)
				}
			}

		case key.Matches(msg, b.keyMap.ShowDirectoriesOnly):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.showDirectoriesOnly = !b.showDirectoriesOnly
				b.showFilesOnly = false

				if b.showDirectoriesOnly {
					return b.getDirectoryListingByTypeCmd(dfs.DirectoriesListingType)
				}

				return b.updateDirectoryListingCmd(dfs.CurrentDirectory)
			}

		case key.Matches(msg, b.keyMap.ShowFilesOnly):
			if b.activeBox == constants.PrimaryBoxActive && !b.showCommandInput && !b.showBoxSpinner {
				b.showFilesOnly = !b.showFilesOnly
				b.showDirectoriesOnly = false

				if b.showFilesOnly {
					return b.getDirectoryListingByTypeCmd(dfs.FilesListingType)
				}

				b.updateDirectoryListingCmd(dfs.CurrentDirectory)
			}

		case key.Matches(msg, b.keyMap.Enter):
			switch {
				case b.findMode:
					b.showCommandInput = false
					b.showBoxSpinner = true

					return b.findFilesByNameCmd(b.textinput.Value())

				case b.receiveMode:
					b.showCommandInput = false
					b.showBoxSpinner = false

					return b.receiveFileCmd(b.textinput.Value())

				default:
					return nil
			}

		case key.Matches(msg, b.keyMap.Edit):
			selectedFile := b.treeFiles[b.treeCursor]

			if !b.showCommandInput && b.activeBox == constants.PrimaryBoxActive && !b.showBoxSpinner {
				if !selectedFile.IsDir() {
					editorPath := cfg.Tran.Editor

					if editorPath == "" {
						return b.handleErrorCmd(errors.New("editor not set, please set it in the config file"))
					}

					editorCmd := exec.Command(editorPath, selectedFile.Name())
					editorCmd.Stdin = os.Stdin
					editorCmd.Stdout = os.Stdout
					editorCmd.Stderr = os.Stderr

					err := editorCmd.Run()
					termenv.AltScreen()

					if err != nil {
						return b.handleErrorCmd(err)
					}

					return tea.Batch(b.redrawCmd(), tea.HideCursor)
				}
			}

		case key.Matches(msg, b.keyMap.Find):
			if !b.showCommandInput && !b.showBoxSpinner {
				b.findMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter a search file/dir name"
				b.textinput.Focus()

				return textinput.Blink
			}

		case key.Matches(msg, b.keyMap.Receive):
			if !b.showCommandInput && !b.showBoxSpinner {
				b.receiveMode = true
				b.showCommandInput = true
				b.textinput.Placeholder = "Enter the password"
				b.textinput.Focus()

				return textinput.Blink
			}

		case key.Matches(msg, b.keyMap.Send):
			b.sendMode = true

		case key.Matches(msg, b.keyMap.Escape):
			b.showCommandInput = false
			b.showFilesOnly = false
			b.showHiddenFiles = true
			b.showDirectoriesOnly = false
			b.findMode = false
			b.sendMode = false
			b.receiveMode = false
			b.errorMsg = ""
			b.foundFilesPaths = nil
			b.showBoxSpinner = false
			b.currentImage = nil
			b.secondaryViewport.GotoTop()
			b.thirdViewport.GotoTop()
			b.textinput.Blur()
			b.textinput.Reset()

			return b.updateDirectoryListingCmd(dfs.CurrentDirectory)

		case key.Matches(msg, b.keyMap.ToggleBox):
			b.activeBox = (b.activeBox + 1) % 2
	}

	b.previousKey = msg

	if b.activeBox != constants.PrimaryBoxActive {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.textinput, cmd = b.textinput.Update(msg)
	cmds = append(cmds, cmd)

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

// handleMouse handles all mouse interaction.
func (b *Bubble) handleMouse(msg tea.MouseMsg) {
	switch msg.Type {
		case tea.MouseWheelUp:
			if b.activeBox == constants.PrimaryBoxActive {
				b.treeCursor--
				b.checkPrimaryViewportBounds()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == constants.SecondaryBoxActive {
				b.secondaryViewport.LineUp(1)
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

		case tea.MouseWheelDown:
			if b.activeBox == constants.PrimaryBoxActive {
				b.treeCursor++
				b.checkPrimaryViewportBounds()
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			if b.activeBox == constants.SecondaryBoxActive {
				b.secondaryViewport.LineDown(1)
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}
		}
}

// Update handles all UI interactions and events for updating the screen.
func (b Bubble) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
		case updateDirectoryListingMsg:
			b.showCommandInput = false
			b.treeCursor = 0
			b.treeFiles = msg
			b.showFileTreePreview = false
			b.fileSizes = make([]string, len(msg))

			for i, file := range msg {
				cmds = append(cmds, b.getDirectoryItemSizeCmd(file.Name(), i))
			}

			b.primaryViewport.SetContent(b.fileTreeView(msg))
			b.textinput.Blur()
			b.textinput.Reset()

			return b, tea.Batch(cmds...)

		case directoryItemSizeMsg:
			if len(b.fileSizes) > 0 && msg.index < len(b.fileSizes) {
				b.fileSizes[msg.index] = msg.size
				b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))
			}

			return b, nil

		case readFileContentMsg:
			b.showFileTreePreview = false
			b.currentImage = nil
			b.secondaryViewport.GotoTop()
			b.thirdViewport.GotoTop()

			switch {
				case msg.code != "":
					b.secondaryBoxContent = msg.code

				case msg.pdfContent != "":
					b.secondaryBoxContent = msg.pdfContent

				case msg.markdown != "":
					b.secondaryBoxContent = msg.markdown

				case msg.image != nil:
					b.currentImage = msg.image
					b.secondaryBoxContent = msg.imageString

				default:
					b.secondaryBoxContent = msg.rawContent
			}

			b.secondaryViewport.SetContent(b.textContentView(b.secondaryBoxContent))

			return b, nil

		case viewDirectoryListingMsg:
			b.showFileTreePreview = true
			b.treePreviewFiles = msg
			b.secondaryViewport.GotoTop()
			b.secondaryViewport.SetContent(b.fileTreePreviewView(msg))

			return b, nil

		case convertImageToStringMsg:
			b.secondaryViewport.GotoTop()
			b.secondaryViewport.SetContent(b.textContentView(string(msg)))

			return b, nil

		case findFilesByNameMsg:
			b.showCommandInput = false
			b.findMode = false
			b.sendMode = false
			b.receiveMode = false
			b.treeCursor = 0
			b.treeFiles = msg.entries
			b.foundFilesPaths = msg.paths
			b.showBoxSpinner = false
			b.textinput.Blur()
			b.textinput.Reset()
			b.fileSizes = make([]string, len(msg.entries))

			for i, file := range msg.entries {
				cmds = append(cmds, b.getDirectoryItemSizeCmd(file.Name(), i))
			}

			b.primaryViewport.SetContent(b.fileTreeView(msg.entries))

			return b, tea.Batch(cmds...)

		case errorMsg:
			b.errorMsg = string(msg)
			b.secondaryViewport.SetContent(b.errorView(string(msg)))

			return b, nil

		case tea.WindowSizeMsg:
			b.width = msg.Width
			b.height = msg.Height

			b.primaryViewport.Width = (msg.Width / 5) - b.primaryViewport.Style.GetHorizontalFrameSize()
			b.primaryViewport.Height = msg.Height - constants.StatusBarHeight - b.primaryViewport.Style.GetVerticalFrameSize()
			b.secondaryViewport.Width = (msg.Width / 2) - b.secondaryViewport.Style.GetHorizontalFrameSize()
			b.secondaryViewport.Height = msg.Height - constants.StatusBarHeight - b.secondaryViewport.Style.GetVerticalFrameSize()
			b.thirdViewport.Width = (msg.Width / 3) - b.thirdViewport.Style.GetHorizontalFrameSize() - 1
			b.thirdViewport.Height = msg.Height - constants.StatusBarHeight - b.thirdViewport.Style.GetVerticalFrameSize()

			b.primaryViewport.SetContent(b.fileTreeView(b.treeFiles))

			hc, _ := glamour.Render(constants.HelpContent, "dark")
			ic, _ := glamour.Render(constants.InfoContent, "dark")

			hs := lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderBottom(true)

			b.thirdViewport.SetContent(hs.Render(hc) + "\n\n" + ic)

			switch {
				case b.showFileTreePreview:
					b.secondaryViewport.SetContent(b.fileTreePreviewView(b.treePreviewFiles))

				case b.currentImage != nil:
					return b, b.convertImageToStringCmd(b.secondaryViewport.Width)

				case b.errorMsg != "":
					b.secondaryViewport.SetContent(b.errorView(b.errorMsg))

				default:
					b.secondaryViewport.SetContent(b.textContentView(b.secondaryBoxContent))
			}

			if !b.ready {
				b.ready = true
			}

			return b, nil

		case tea.MouseMsg:
			b.handleMouse(msg)

		case tea.KeyMsg:
			cmd = b.handleKeys(msg)
			cmds = append(cmds, cmd)

			return b, tea.Batch(cmds...)
	}

	if b.activeBox != constants.PrimaryBoxActive {
		b.secondaryViewport, cmd = b.secondaryViewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	b.textinput, cmd = b.textinput.Update(msg)
	cmds = append(cmds, cmd)

	b.spinner, cmd = b.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return b, tea.Batch(cmds...)
}

package tui

import (
	"os"
	"image"
	"io/fs"
	_ "image/png"
	_ "image/jpeg"
	"path/filepath"

	"github.com/abdfnx/tran/dfs"
	"github.com/abdfnx/tran/models"
	"github.com/abdfnx/tran/renderer"
	"github.com/abdfnx/tran/constants"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
)

type updateDirectoryListingMsg []fs.DirEntry
type viewDirectoryListingMsg []fs.DirEntry
type errorMsg string
type convertImageToStringMsg string
type directoryItemSizeMsg struct {
	index int
	size  string
}

type findFilesByNameMsg struct {
	paths   []string
	entries []fs.DirEntry
}

type readFileContentMsg struct {
	rawContent  string
	markdown    string
	code        string
	imageString string
	pdfContent  string
	image       image.Image
}

// updateDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (b Bubble) updateDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		files, err := dfs.GetDirectoryListing(name, b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		err = os.Chdir(name)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(files)
	}
}

// viewDirectoryListingCmd updates the directory listing based on the name of the directory provided.
func (b Bubble) viewDirectoryListingCmd(name string) tea.Cmd {
	return func() tea.Msg {
		currentDir, err := dfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		files, err := dfs.GetDirectoryListing(filepath.Join(currentDir, name), b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		return viewDirectoryListingMsg(files)
	}
}

// convertImageToStringCmd redraws the image based on the width provided.
func (b Bubble) convertImageToStringCmd(width int) tea.Cmd {
	return func() tea.Msg {
		imageString := renderer.ImageToString(width, b.currentImage)

		return convertImageToStringMsg(imageString)
	}
}

// readFileContentCmd reads the content of a file and returns it.
func (b Bubble) readFileContentCmd(fileName string, width int) tea.Cmd {
	return func() tea.Msg {
		content, err := dfs.ReadFileContent(fileName)

		if err != nil {
			return errorMsg(err.Error())
		}

		switch {
			case filepath.Ext(fileName) == ".md":
				markdownContent, err := renderer.RenderMarkdown(width, content)

				if err != nil {
					return errorMsg(err.Error())
				}

				return readFileContentMsg{
					rawContent:  content,
					markdown:    markdownContent,
					code:        "",
					imageString: "",
					pdfContent:  "",
					image:       nil,
				}

			case filepath.Ext(fileName) == ".png" || filepath.Ext(fileName) == ".jpg" || filepath.Ext(fileName) == ".jpeg":
				imageContent, err := os.Open(fileName)

				if err != nil {
					return errorMsg(err.Error())
				}

				img, _, err := image.Decode(imageContent)

				if err != nil {
					return errorMsg(err.Error())
				}

				imageString := renderer.ImageToString(width, img)

				return readFileContentMsg{
					rawContent:  content,
					code:        "",
					markdown:    "",
					imageString: imageString,
					pdfContent:  "",
					image:       img,
				}

			case filepath.Ext(fileName) == ".pdf":
				pdfContent, err := renderer.ReadPdf(fileName)
				if err != nil {
					return errorMsg(err.Error())
				}

				return readFileContentMsg{
					rawContent:  content,
					code:        "",
					markdown:    "",
					imageString: "",
					pdfContent:  pdfContent,
					image:       nil,
				}

			default:
				syntaxTheme := "solarized-light"

				if lipgloss.HasDarkBackground() {
					syntaxTheme = "solarized-dark"
				}

				code, err := renderer.Highlight(content, filepath.Ext(fileName), syntaxTheme)

				if err != nil {
					return errorMsg(err.Error())
				}

				return readFileContentMsg{
					rawContent:  content,
					code:        code,
					markdown:    "",
					imageString: "",
					pdfContent:  "",
					image:       nil,
				}
		}
	}
}

// getDirectoryItemSizeCmd calculates the size of a directory or file.
func (b Bubble) getDirectoryItemSizeCmd(name string, i int) tea.Cmd {
	return func() tea.Msg {
		size, err := dfs.GetDirectoryItemSize(name)

		if err != nil {
			return directoryItemSizeMsg{size: "N/A", index: i}
		}

		sizeString := renderer.ConvertBytesToSizeString(size)

		return directoryItemSizeMsg{
			size:  sizeString,
			index: i,
		}
	}
}

// handleErrorCmd returns an error message to the UI.
func (b Bubble) handleErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err.Error())
	}
}

// getDirectoryListingByTypeCmd returns only directories in the current directory.
func (b Bubble) getDirectoryListingByTypeCmd(listType string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		directories, err := dfs.GetDirectoryListingByType(workingDir, listType, b.showHiddenFiles)
		if err != nil {
			return errorMsg(err.Error())
		}

		return updateDirectoryListingMsg(directories)
	}
}

// findFilesByNameCmd finds files based on name.
func (b Bubble) findFilesByNameCmd(name string) tea.Cmd {
	return func() tea.Msg {
		workingDir, err := dfs.GetWorkingDirectory()
		if err != nil {
			return errorMsg(err.Error())
		}

		paths, entries, err := dfs.FindFilesByName(name, workingDir)
		if err != nil {
			return errorMsg(err.Error())
		}

		return findFilesByNameMsg{
			paths:   paths,
			entries: entries,
		}
	}
}

func (b Bubble) receiveFileCmd(password string) tea.Cmd {
	return func() tea.Msg {
		err := ValidateTranxAddress()

		if err != nil {
			return err
		}

		HandleReceiveCommand(models.TranOptions{
			TranxAddress: constants.DEFAULT_ADDRESS,
			TranxPort:    constants.DEFAULT_PORT,
		}, password)

		return nil
	}
}

// redrawCmd redraws the UI.
func (b Bubble) redrawCmd() tea.Cmd {
	return func() tea.Msg {
		return tea.WindowSizeMsg{
			Width:  b.width,
			Height: b.height,
		}
	}
}

package tui

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/abdfnx/tran/dfs"
	"github.com/abdfnx/tran/tools"
	"github.com/abdfnx/tran/models"
	"github.com/abdfnx/tran/renderer"
	"github.com/abdfnx/tran/constants"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/truncate"
)

// statusBarView returns the status bar.
func (b Bubble) statusBarView() string {
	var logo string
	var status string

	width := lipgloss.Width
	selectedFileName := "N/A"
	fileCount := "0/0"

	if len(b.treeFiles) > 0 && b.treeFiles[b.treeCursor] != nil {
		selectedFile, err := b.treeFiles[b.treeCursor].Info()
		if err != nil {
			return "error"
		}

		fileCount = fmt.Sprintf("%d/%d", b.treeCursor+1, len(b.treeFiles))
		selectedFileName = selectedFile.Name()

		currentPath, err := dfs.GetWorkingDirectory()
		if err != nil {
			currentPath = dfs.CurrentDirectory
		}

		if len(b.foundFilesPaths) > 0 {
			currentPath = b.foundFilesPaths[b.treeCursor]
		}

		status = fmt.Sprintf("%s %s",
			selectedFile.ModTime().Format("2006-01-02 15:04:05"),
			currentPath,
		)
	}

	if b.showCommandInput {
		status = b.textinput.View()
	}

	logo = "TRAN"

	// Selected file styles
	selectedFileStyle := constants.BoldTextStyle.Copy().
		Foreground(b.theme.StatusBarSelectedFileForegroundColor).
		Background(b.theme.StatusBarSelectedFileBackgroundColor)

	selectedFileColumn := selectedFileStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(truncate.StringWithTail(selectedFileName, 30, "..."))

	// File count styles
	fileCountStyle := constants.BoldTextStyle.Copy().
		Foreground(b.theme.StatusBarTotalFilesForegroundColor).
		Background(b.theme.StatusBarTotalFilesBackgroundColor)

	fileCountColumn := fileCountStyle.
		Align(lipgloss.Right).
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(fileCount)

	// Logo styles
	logoStyle := constants.BoldTextStyle.Copy().
		Foreground(b.theme.StatusBarLogoForegroundColor).
		Background(b.theme.StatusBarLogoBackgroundColor)

	logoColumn := logoStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Render(logo)

	// Status styles
	statusStyle := constants.BoldTextStyle.Copy().
		Foreground(b.theme.StatusBarBarForegroundColor).
		Background(b.theme.StatusBarBarBackgroundColor)

	statusColumn := statusStyle.
		Padding(0, 1).
		Height(constants.StatusBarHeight).
		Width(b.width - width(selectedFileColumn) - width(fileCountColumn) - width(logoColumn)).
		Render(truncate.StringWithTail(
			status,
			uint(b.width-width(selectedFileColumn)-width(fileCountColumn)-width(logoColumn)-3),
			"..."),
		)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		logoColumn,
		selectedFileColumn,
		statusColumn,
		fileCountColumn,
	)
}

// fileView returns the filetree view.
func (b Bubble) fileTreeView(files []fs.DirEntry) string {
	var directoryItem string
	curFiles := ""
	fileSize := ""
	selectedItemColor := b.theme.SelectedTreeItemColor
	unselectedItemColor := b.theme.UnselectedTreeItemColor

	for i, file := range files {
		fileInfo, err := file.Info()

		if err != nil {
			return "Error loading directory tree"
		}

		if b.treeCursor == i {
			if len(b.fileSizes) > 0 {
				if b.fileSizes[i] != "" {
					fileSize = constants.BoldTextStyle.Copy().
						Foreground(constants.Colors["black"]).
						Background(selectedItemColor).
						Render(b.fileSizes[i])
				} else {
					fileSize = constants.BoldTextStyle.Copy().
						Foreground(constants.Colors["black"]).
						Background(selectedItemColor).
						Render(constants.FileSizeLoadingStyle)
				}
			}

			directoryItem = constants.BoldTextStyle.Copy().
				Background(selectedItemColor).
				Width(b.primaryViewport.Width - lipgloss.Width(fileSize)).
				Foreground(constants.Colors["black"]).
				Render(
					truncate.StringWithTail(
						fileInfo.Name(), uint(b.primaryViewport.Width-lipgloss.Width(fileSize)), constants.EllipsisStyle,
					),
				)
		} else {
			if len(b.fileSizes) > 0 {
				if b.fileSizes[i] != "" {
					fileSize = constants.BoldTextStyle.Copy().
						Foreground(unselectedItemColor).
						Render(b.fileSizes[i])
				} else {
					fileSize = constants.BoldTextStyle.Copy().
						Foreground(unselectedItemColor).
						Render(constants.FileSizeLoadingStyle)
				}
			}

			directoryItem = constants.BoldTextStyle.Copy().
				Width(b.primaryViewport.Width - lipgloss.Width(fileSize)).
				Foreground(unselectedItemColor).
				Render(
					truncate.StringWithTail(
						fileInfo.Name(), uint(b.primaryViewport.Width-lipgloss.Width(fileSize)), constants.EllipsisStyle,
					),
				)
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, directoryItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	if len(files) == 0 {
		curFiles = "Directory is empty"
	}

	return lipgloss.NewStyle().
		Width(b.primaryViewport.Width).
		Height(b.primaryViewport.Height).
		Render(curFiles)
}

// fileTreePreviewView returns a preview of a filetree.
func (b Bubble) fileTreePreviewView(files []fs.DirEntry) string {
	var directoryItem string
	curFiles := ""

	for _, file := range files {
		fileColor := b.theme.UnselectedTreeItemColor

		fileInfo, _ := file.Info()

		fileSize := lipgloss.NewStyle().
			Foreground(fileColor).
			Render(renderer.ConvertBytesToSizeString(fileInfo.Size()))

		directoryItem = constants.BoldTextStyle.Copy().
				Foreground(fileColor).
				Render(fileInfo.Name())

		dirItem := lipgloss.NewStyle().Width(
			b.secondaryViewport.Width - lipgloss.Width(fileSize),
		).Render(
			truncate.StringWithTail(
				directoryItem, uint(b.secondaryViewport.Width-lipgloss.Width(fileSize)), constants.EllipsisStyle,
			),
		)

		row := lipgloss.JoinHorizontal(lipgloss.Top, dirItem, fileSize)

		curFiles += fmt.Sprintf("%s\n", row)
	}

	if len(files) == 0 {
		curFiles = "Directory is empty"
	}

	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(curFiles)
}

// textContentView returns some text content.
func (b Bubble) textContentView(content string) string {
	return lipgloss.NewStyle().
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(content)
}

// errorView returns an error message.
func (b Bubble) errorView(msg string) string {
	return lipgloss.NewStyle().
		Foreground(b.theme.ErrorColor).
		Width(b.secondaryViewport.Width).
		Height(b.secondaryViewport.Height).
		Render(msg)
}

// View returns a string representation of the entire application UI.
func (b Bubble) View() string {
	if !b.ready {
		return fmt.Sprintf("%s %s", b.spinner.View(), "loading...")
	}

	if b.sendMode {
		selectedFile := b.treeFiles[b.treeCursor]

		tools.RandomSeed()

		err := ValidateTranxAddress()

		if err != nil {
			fmt.Println(err)
		}

		fn := []string{selectedFile.Name()}

		HandleSendCommand(models.TranOptions{
			TranxAddress: constants.DEFAULT_ADDRESS,
			TranxPort:    constants.DEFAULT_PORT,
		}, fn)
	}

	var primaryBox string
	var secondaryBox string
	var thirdBox string

	primaryBoxBorder := lipgloss.RoundedBorder()
	secondaryBoxBorder := lipgloss.RoundedBorder()
	thirdBoxBorder := lipgloss.RoundedBorder()
	primaryBoxBorderColor := b.theme.InactiveBoxBorderColor
	secondaryBoxBorderColor := b.theme.InactiveBoxBorderColor
	thirdBoxBorderColor := b.theme.InactiveBoxBorderColor

	if b.activeBox == constants.PrimaryBoxActive {
		primaryBoxBorderColor = b.theme.ActiveBoxBorderColor
	}

	if b.activeBox == constants.SecondaryBoxActive {
		secondaryBoxBorderColor = b.theme.ActiveBoxBorderColor
	}

	if b.activeBox == constants.ThirdBoxActive {
		thirdBoxBorderColor = b.theme.ActiveBoxBorderColor
	}

	if b.appConfig.Tran.Borderless {
		primaryBoxBorder = lipgloss.HiddenBorder()
		secondaryBoxBorder = lipgloss.HiddenBorder()
		thirdBoxBorder = lipgloss.HiddenBorder()
	}

	b.primaryViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder).
		BorderForeground(primaryBoxBorderColor)

	primaryBox = b.primaryViewport.View()

	if b.showBoxSpinner {
		b.primaryViewport.Style = lipgloss.NewStyle().
			PaddingLeft(constants.BoxPadding).
			PaddingRight(constants.BoxPadding).
			Border(primaryBoxBorder).
			BorderForeground(primaryBoxBorderColor)

		primaryBox = b.primaryViewport.Style.Render(fmt.Sprintf("%s loading...", b.spinner.View()))
	}

	b.secondaryViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder).
		BorderForeground(secondaryBoxBorderColor)

	secondaryBox = b.secondaryViewport.View()

	b.thirdViewport.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(thirdBoxBorder).
		BorderForeground(thirdBoxBorderColor)
	
	thirdBox = b.thirdViewport.View()

	s := strings.Builder{}

	view := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			primaryBox,
			secondaryBox,
			thirdBox,
		),
	)

	s.WriteString(view)
	s.WriteString("\n")
	s.WriteString(b.statusBarView())

	return s.String()
}

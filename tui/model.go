package tui

import (
	"image"
	"io/fs"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/scmn-dev/tran/config"
	"github.com/scmn-dev/tran/constants"
	"github.com/scmn-dev/tran/theme"
)

// Bubble represents the state of the UI.
type Bubble struct {
	appConfig              config.Config
	theme                  theme.Theme
	currentImage           image.Image
	spinner                spinner.Model
	textinput              textinput.Model
	primaryViewport        viewport.Model
	secondaryViewport      viewport.Model
	thirdViewport          viewport.Model
	treeFiles              []fs.DirEntry
	treePreviewFiles       []fs.DirEntry
	previousKey            tea.KeyMsg
	keyMap                 KeyMap
	width                  int
	height                 int
	activeBox              int
	treeCursor             int
	showHiddenFiles        bool
	ready                  bool
	showCommandInput       bool
	showFilesOnly          bool
	showDirectoriesOnly    bool
	showFileTreePreview    bool
	findMode               bool
	sendMode			   bool
	receiveMode			   bool
	showBoxSpinner         bool
	foundFilesPaths        []string
	fileSizes              []string
	secondaryBoxContent    string
	errorMsg               string
}

// New creates an instance of the entire application.
func New() Bubble {
	cfg := config.GetConfig()
	theme := theme.GetTheme("default")

	primaryBoxBorder := lipgloss.RoundedBorder()
	secondaryBoxBorder := lipgloss.RoundedBorder()
	thirdBoxBorder := lipgloss.RoundedBorder()
	primaryBoxBorderColor := theme.ActiveBoxBorderColor
	secondaryBoxBorderColor := theme.InactiveBoxBorderColor
	thirdBoxBorderColor := theme.InactiveBoxBorderColor

	if cfg.Tran.Borderless {
		primaryBoxBorder = lipgloss.HiddenBorder()
		secondaryBoxBorder = lipgloss.HiddenBorder()
		thirdBoxBorder = lipgloss.HiddenBorder()
	}

	pvp := viewport.New(0, 0)
	pvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(primaryBoxBorder).
		BorderForeground(primaryBoxBorderColor)

	svp := viewport.New(0, 0)
	svp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(secondaryBoxBorder).
		BorderForeground(secondaryBoxBorderColor)

	tvp := viewport.New(0, 0)
	tvp.Style = lipgloss.NewStyle().
		PaddingLeft(constants.BoxPadding).
		PaddingRight(constants.BoxPadding).
		Border(thirdBoxBorder).
		BorderForeground(thirdBoxBorderColor)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(theme.SpinnerColor)

	t := textinput.New()
	t.Prompt = "❯ "
	t.CharLimit = 250
	t.PlaceholderStyle = lipgloss.NewStyle().
		Background(theme.StatusBarBarBackgroundColor).
		Foreground(theme.StatusBarBarForegroundColor)

	return Bubble{
		appConfig:         cfg,
		theme:             theme,
		showHiddenFiles:   true,
		spinner:           s,
		textinput:         t,
		primaryViewport:   pvp,
		secondaryViewport: svp,
		thirdViewport:     tvp,
		keyMap:            Keys,
	}
}

package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/abdfnx/tran/constants"
	"github.com/abdfnx/tran/tools"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
)

// ui state flows from the top down
const (
	showPasswordWithCopy uiState = iota
	showPassword
	showSendingProgress
	showSFinished
	showSError
)

type senderUIModel struct {
	state        uiState
	fileNames    []string
	payloadSize  int64
	password     string
	readyToSend  bool
	spinner      spinner.Model
	progressBar  progress.Model
	errorMessage string
}

type ReadyMsg struct{}

type PasswordMsg struct {
	Password string
}

func NewSenderUI() *tea.Program {
	m := senderUIModel{progressBar: constants.ProgressBar}
	m.resetSpinner()
	var opts []tea.ProgramOption

	termenv.AltScreen()

	opts = append(opts, tea.WithAltScreen())

	return tea.NewProgram(m, opts...)
}

func (senderUIModel) Init() tea.Cmd {
	return spinner.Tick
}

func (m senderUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case FileInfoMsg:
			m.fileNames = msg.FileNames
			m.payloadSize = msg.Bytes

			return m, nil

		case ReadyMsg:
			m.readyToSend = true
			m.resetSpinner()

			return m, spinner.Tick

		case PasswordMsg:
			m.password = msg.Password

			return m, nil

		case ProgressMsg:
			if m.state != showSendingProgress {
				m.state = showSendingProgress
				m.resetSpinner()
				return m, spinner.Tick
			}

			if m.progressBar.Percent() == 1.0 {
				return m, nil
			}

			cmd := m.progressBar.SetPercent(float64(msg.Progress))

			return m, cmd

		case FinishedMsg:
			m.state = showSFinished
			cmd := m.progressBar.SetPercent(1.0)

			return m, cmd

		case ErrorMsg:
			m.state = showSError
			m.errorMessage = msg.Message

			return m, nil

		case tea.KeyMsg:
			if tools.Contains(constants.QuitKeys, strings.ToLower(msg.String())) {
				return m, tea.Quit
			}

			return m, nil

		case tea.WindowSizeMsg:
			m.progressBar.Width = msg.Width - 2 * constants.PADDING - 4

			if m.progressBar.Width > constants.MAX_WIDTH {
				m.progressBar.Width = constants.MAX_WIDTH
			}

			return m, nil

		// FrameMsg is sent when the progress bar wants to animate itself
		case progress.FrameMsg:
			progressModel, cmd := m.progressBar.Update(msg)
			m.progressBar = progressModel.(progress.Model)

			return m, cmd

		default:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)

			return m, cmd
	}
}

func (m senderUIModel) View() string {
	readiness := fmt.Sprintf("%s Compressing objects, preparing to send", m.spinner.View())

	if m.readyToSend {
		readiness = fmt.Sprintf("%s Awaiting receiver, ready to send", m.spinner.View())
	}

	if m.state == showSendingProgress {
		readiness = fmt.Sprintf("%s Sending", m.spinner.View())
	}

	fileInfoText := fmt.Sprintf("%s object(s)...", readiness)

	if m.fileNames != nil && m.payloadSize != 0 {
		sort.Strings(m.fileNames)
		filesToSend := constants.ItalicText(strings.Join(m.fileNames, ", "))
		payloadSize := constants.BoldText(tools.ByteCountSI(m.payloadSize))
		fileInfoText = fmt.Sprintf("%s %d objects (%s)", readiness, len(m.fileNames), payloadSize)

		indentedWrappedFiles := indent.String(wordwrap.String(fmt.Sprintf("Sending: %s", filesToSend), constants.MAX_WIDTH), constants.PADDING)
		fileInfoText = fmt.Sprintf("%s\n\n%s", fileInfoText, indentedWrappedFiles)
	}

	switch m.state {
		case showPassword, showPasswordWithCopy:
			return "\n" +
				constants.PadText + constants.InfoStyle(fileInfoText) + "\n\n" +
				constants.PadText + "On the other computer, press " + constants.HelpStyle("`ctrl+r`") + " to enable receive mode and then enter the password:" + "\n\n" +
				constants.PadText + "This is the passowrd: " + constants.BoldText(m.password) + "\n\n"

		case showSendingProgress:
			return "\n" +
				constants.PadText + constants.InfoStyle(fileInfoText) + "\n\n" +
				constants.PadText + m.progressBar.View() + "\n\n" +
				constants.PadText + constants.QuitCommandsHelpText + "\n\n"

		case showSFinished:
			payloadSize := constants.BoldText(tools.ByteCountSI(m.payloadSize))
			indentedWrappedFiles := indent.String(fmt.Sprintf("Sent: %s", wordwrap.String(constants.ItalicText(TopLevelFilesText(m.fileNames)), constants.MAX_WIDTH)), constants.PADDING)
			finishedText := fmt.Sprintf("Sent %d objects (%s decompressed)\n\n%s", len(m.fileNames), payloadSize, indentedWrappedFiles)

			return "\n" +
				constants.PadText + constants.InfoStyle(finishedText) + "\n\n" +
				constants.PadText + m.progressBar.View() + "\n\n" +
				constants.PadText + constants.QuitCommandsHelpText + "\n\n"

		case showSError:
			return m.errorMessage

		default:
			return ""
	}
}

func (m *senderUIModel) resetSpinner() {
	m.spinner = spinner.NewModel()
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(constants.PRIMARY_COLOR))

	if m.readyToSend {
		m.spinner.Spinner = WaitingSpinner
	} else {
		m.spinner.Spinner = CompressingSpinner
	}

	if m.state == showSendingProgress {
		m.spinner.Spinner = TransferSpinner
	}
}

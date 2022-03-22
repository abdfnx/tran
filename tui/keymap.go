package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit                key.Binding
	Down                key.Binding
	Up                  key.Binding
	Left                key.Binding
	Right               key.Binding
	View                key.Binding
	Receive			    key.Binding
	GotoBottom          key.Binding
	HomeShortcut        key.Binding
	RootShortcut        key.Binding
	ToggleHidden        key.Binding
	ShowDirectoriesOnly key.Binding
	ShowFilesOnly       key.Binding
	Enter               key.Binding
	Edit                key.Binding
	Find                key.Binding
	Send 			    key.Binding
	Command             key.Binding
	Escape              key.Binding
	ToggleBox           key.Binding
}

// DefaultKeyMap returns a set of default keybindings.
var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
	),
	Up: key.NewBinding(
		key.WithKeys("up"),
	),
	Left: key.NewBinding(
		key.WithKeys("left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
	),
	View: key.NewBinding(
		key.WithKeys("V"),
	),
	GotoBottom: key.NewBinding(
		key.WithKeys("G"),
	),
	HomeShortcut: key.NewBinding(
		key.WithKeys("~"),
	),
	RootShortcut: key.NewBinding(
		key.WithKeys("/"),
	),
	ToggleHidden: key.NewBinding(
		key.WithKeys("."),
	),
	ShowDirectoriesOnly: key.NewBinding(
		key.WithKeys("D"),
	),
	ShowFilesOnly: key.NewBinding(
		key.WithKeys("F"),
	),
	Edit: key.NewBinding(
		key.WithKeys("E"),
	),
	Find: key.NewBinding(
		key.WithKeys("ctrl+f"),
	),
	ToggleBox: key.NewBinding(
		key.WithKeys("tab"),
	),
	Receive: key.NewBinding(
		key.WithKeys("ctrl+r"),
	),
	Send: key.NewBinding(
		key.WithKeys("ctrl+s"),
	),
}

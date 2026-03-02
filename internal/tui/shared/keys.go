package shared

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings for the application.
type KeyMap struct {
	Quit         key.Binding
	Help         key.Binding
	Tab          key.Binding
	ShiftTab     key.Binding
	Enter        key.Binding
	Escape       key.Binding
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	Home         key.Binding
	End          key.Binding
	NextPage     key.Binding
	PrevPage     key.Binding
	ToggleEditor key.Binding
	ExecSQL      key.Binding
	EditCell     key.Binding
	InsertRow    key.Binding
	DeleteRow    key.Binding
}

// DefaultKeyMap returns the default key bindings.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next pane"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev pane"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Escape: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup"),
			key.WithHelp("pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown"),
			key.WithHelp("pgdown", "page down"),
		),
		Home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "first column"),
		),
		End: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "last column"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "next page"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "prev page"),
		),
		ToggleEditor: key.NewBinding(
			key.WithKeys("ctrl+e"),
			key.WithHelp("ctrl+e", "SQL editor"),
		),
		ExecSQL: key.NewBinding(
			key.WithKeys("alt+enter", "f5"),
			key.WithHelp("alt+enter/F5", "execute SQL"),
		),
		EditCell: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit cell"),
		),
		InsertRow: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert row"),
		),
		DeleteRow: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete row"),
		),
	}
}

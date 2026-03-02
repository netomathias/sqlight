package shared

import "github.com/charmbracelet/lipgloss"

// Styles holds all application styles.
type Styles struct {
	App       lipgloss.Style
	Sidebar   SidebarStyles
	Table     TableStyles
	StatusBar StatusBarStyles
	Editor    EditorStyles
	Dialog    DialogStyles
	Help      HelpStyles
}

type SidebarStyles struct {
	Border       lipgloss.Style
	Title        lipgloss.Style
	Item         lipgloss.Style
	SelectedItem lipgloss.Style
	TypeBadge    lipgloss.Style
}

type TableStyles struct {
	Border   lipgloss.Style
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
	Null     lipgloss.Style
}

type EditorStyles struct {
	Border lipgloss.Style
	Title  lipgloss.Style
}

type StatusBarStyles struct {
	Bar     lipgloss.Style
	Key     lipgloss.Style
	Value   lipgloss.Style
	Error   lipgloss.Style
	Success lipgloss.Style
}

type DialogStyles struct {
	Border  lipgloss.Style
	Title   lipgloss.Style
	Input   lipgloss.Style
	Button  lipgloss.Style
	Overlay lipgloss.Style
}

type HelpStyles struct {
	Key  lipgloss.Style
	Desc lipgloss.Style
}

// DefaultStyles returns the default application styles.
func DefaultStyles() Styles {
	accent := lipgloss.Color("4")       // blue
	dimmed := lipgloss.Color("8")       // gray
	highlight := lipgloss.Color("11")   // yellow
	errorColor := lipgloss.Color("1")   // red
	successColor := lipgloss.Color("2") // green

	return Styles{
		App: lipgloss.NewStyle(),
		Sidebar: SidebarStyles{
			Border: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(dimmed).
				Padding(0, 1),
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent).
				Padding(0, 1),
			Item: lipgloss.NewStyle().
				Padding(0, 1),
			SelectedItem: lipgloss.NewStyle().
				Padding(0, 1).
				Bold(true).
				Foreground(highlight).
				Background(lipgloss.Color("0")),
			TypeBadge: lipgloss.NewStyle().
				Foreground(dimmed),
		},
		Table: TableStyles{
			Border: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(dimmed),
			Header: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent).
				Padding(0, 1),
			Cell: lipgloss.NewStyle().
				Padding(0, 1),
			Selected: lipgloss.NewStyle().
				Padding(0, 1).
				Bold(true).
				Background(lipgloss.Color("0")),
			Null: lipgloss.NewStyle().
				Foreground(dimmed).
				Italic(true),
		},
		Editor: EditorStyles{
			Border: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(dimmed).
				Padding(0, 1),
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent),
		},
		StatusBar: StatusBarStyles{
			Bar: lipgloss.NewStyle().
				Padding(0, 1).
				Background(lipgloss.Color("0")),
			Key: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent),
			Value: lipgloss.NewStyle().
				Foreground(lipgloss.Color("7")),
			Error: lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true),
			Success: lipgloss.NewStyle().
				Foreground(successColor),
		},
		Dialog: DialogStyles{
			Border: lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(accent).
				Padding(1, 2),
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent),
			Input: lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(dimmed).
				Padding(0, 1),
			Button: lipgloss.NewStyle().
				Padding(0, 2).
				Background(accent).
				Foreground(lipgloss.Color("0")),
			Overlay: lipgloss.NewStyle(),
		},
		Help: HelpStyles{
			Key: lipgloss.NewStyle().
				Bold(true).
				Foreground(accent),
			Desc: lipgloss.NewStyle().
				Foreground(dimmed),
		},
	}
}

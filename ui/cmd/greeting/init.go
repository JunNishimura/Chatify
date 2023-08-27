package greeting

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/gdamore/tcell/v2" // if not import this package, UI does not display as expected
)

func (m *Model) Init() tea.Cmd {
	return tea.Batch(m.loadConfig, textinput.Blink)
}

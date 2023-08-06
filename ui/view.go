package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	doneStyle = lipgloss.NewStyle().Margin(1, 2)
)

func (m *Model) View() string {
	view := m.viewport.View()
	text := m.textarea.View()

	var output string
	if m.greetingDone {
		output += doneStyle.Render(fmt.Sprintf(`Nice to see you, %s!!
	
If you want to talk to me, please type

  $ chatify hey

See you then!! (press to "enter" to exit)`,
			m.user.DisplayName,
		))
	} else if m.qaDone {
		for i, key := range m.confKeyList[:m.writeIndex] {
			output += fmt.Sprintf("%s set %s as %s\n", checkMark, key, m.qaList[i].answer)
		}

		if m.setConfigDone {
			output += doneStyle.Render(`Press "enter" to authorize`)
		}
	} else {
		output = fmt.Sprintf("%s\n\n%s\n\n", view, text)
	}

	return output
}

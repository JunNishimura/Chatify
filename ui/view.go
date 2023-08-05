package ui

import "fmt"

func (m *Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n",
		m.viewport.View(),
		m.textarea.View(),
	)
}

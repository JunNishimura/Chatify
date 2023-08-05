package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			m.qaList[m.index].answer = m.textarea.Value()

			m.index++
			if m.index == len(m.qaList) {
				m.displayMessages = append(m.displayMessages, m.senderStyle.Render("Chatify: ")+"Thank you so much!")
				m.viewport.SetContent(strings.Join(m.displayMessages, "\n"))

				return m, tea.Quit
			}

			m.displayMessages = append(
				m.displayMessages,
				m.senderStyle.Render("You: ")+m.textarea.Value(),
				m.senderStyle.Render("Chatify: ")+m.qaList[m.index].question,
			)
			m.viewport.SetContent(strings.Join(m.displayMessages, "\n"))
			m.viewport.GotoBottom()

			m.textarea.Reset()
			m.textarea.Placeholder = m.qaList[m.index].placeholder
		}
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

package program

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type md struct {
	textInput textinput.Model
	err       error
}

func InitialModel() md {
	ti := textinput.New()
	ti.Placeholder = "Your Message?"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return md{
		textInput: ti,
		err:       nil,
	}
}

func (m md) Init() tea.Cmd {
	return textinput.Blink
}

func (m md) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			fmt.Println(m.textInput.Value())
			return m, tea.Println("done")
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m md) View() string {
	return fmt.Sprintf(
		"\n%s\n\n%s",
		m.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

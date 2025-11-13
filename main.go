package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/notepad_pro/db"
	repo "github.com/notepad_pro/model"

	"github.com/charmbracelet/lipgloss"
	"os"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("229"))

	rowStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("252"))

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1)
)

type model struct {
	list []*repo.Todo
}

func (m model) Init() tea.Cmd {
	return nil
}

func renderList(list []*repo.Todo) string {
	if len(list) == 0 {
		return "No todos yet!"
	}

	var rows []string

	header := headerStyle.Render(fmt.Sprintf("%-4s %-30s %-20s", "ID", "Task", "Created At"))
	rows = append(rows, header)

	for _, t := range list {
		row := rowStyle.Render(fmt.Sprintf("%-4d %-30s %-20s",
			t.ID,
			t.Task,
			t.CreatedAt.Format("2006-01-02 15:04"),
		))
		rows = append(rows, row)
	}

	table := lipgloss.JoinVertical(lipgloss.Left, rows...)

	return borderStyle.Render(table)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return renderList(m.list)
}

func main() {
	defer db.DB.Close()

	todos, _ := repo.GetList(1)
	p := tea.NewProgram(model{
		list: todos,
	})

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

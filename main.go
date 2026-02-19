package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/notepad_pro/db"
	repo "github.com/notepad_pro/model"
)

// headerStyle = lipgloss.NewStyle().
// 		Bold(true).
// 		Foreground(lipgloss.Color("229"))
//
// rowStyle = lipgloss.NewStyle().
// 		Padding(0, 1).
// 		Foreground(lipgloss.Color("252"))

var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(1)

type Model struct {
	Table     table.Model
	TextInput textinput.Model
	ShowInput bool
	Height    int
	Width     int
	Message   string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.ShowInput {
				id := len(m.Table.Rows())
				task := m.TextInput.Value()
				errVar := appendRow(&repo.Todo{
					ID:        uint(id),
					Task:      task,
					CreatedAt: time.Now(),
				})
				if errVar != nil {
					m.Message = fmt.Sprintf("%v", errVar)
				}
				m := renderTable(1, false)
				return m, nil
			} else {
				return m, tea.Batch(
					tea.Printf("Resume task : %s", m.Table.SelectedRow()[1]),
				)
			}
		case "ctrl+d":
			id, err := strconv.Atoi(m.Table.SelectedRow()[0])
			if err != nil {
				log.Fatalf("error in convertion : %v", err)
			}
			err = deleteRow(id)
			if err != nil {
				log.Fatalf("error deleting row from table : %v", err)
			}
			idx := m.Table.Cursor()
			rows := m.Table.Rows()
			if len(rows) == 0 {
				return m, nil
			}
			rows = append(rows[:idx], rows[idx+1:]...)
			m.Table.SetRows(rows)
		case "ctrl+a":
			m := renderTable(1, true)
			return m, cmd
		}
	}
	if m.ShowInput {
		m.TextInput, cmd = m.TextInput.Update(msg)
	} else {
		m.Table, cmd = m.Table.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	s := m.Table.View()
	if m.ShowInput {
		s += "\n\n" + m.TextInput.View()
	}

	tableString := borderStyle.Render(s)

	margin := 2

	box := lipgloss.NewStyle().
		Width(m.Width - margin*2).
		Height(m.Height - margin*2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("46")).
		Render(tableString)

	liststring, _ := repo.GetRawList()
	newJson, _ := json.Marshal(liststring)

	return lipgloss.Place(
		m.Width,
		m.Height,
		lipgloss.Center,
		lipgloss.Center,
		box+"\n"+string(newJson)+m.Message,
	)
}

func deleteRow(id int) error {
	err := repo.DeleteTodo(id)
	if err != nil {
		return err
	}
	return nil
}

func appendRow(todo *repo.Todo) error {
	err := repo.CreateTodo(todo)
	if err != nil {
		return err
	}
	return nil
}

func renderTable(page int, showInput bool) Model {
	todos, _ := repo.GetRawList()

	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Task", Width: 30},
		{Title: "Created At", Width: 10},
	}

	rows := []table.Row{}

	for _, t := range todos {
		rows = append(rows, table.Row{
			strconv.Itoa(int(t.ID)),
			t.Task,
			t.CreatedAt.Format(time.RFC822),
		})
	}

	contentTable := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	contentTable.SetStyles(s)
	m := Model{}
	if showInput {
		ti := textinput.New()
		ti.Placeholder = "enter a task..."
		ti.Focus()
		ti.CharLimit = 156
		ti.Width = 20
		m.TextInput = ti
		contentTable.Blur()
		m.Table = contentTable
		m.ShowInput = showInput
	} else {
		m.Table = contentTable
		m.ShowInput = showInput
	}
	return m
}

func main() {
	defer db.DB.Close()

	t := renderTable(1, false)
	p := tea.NewProgram(t, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

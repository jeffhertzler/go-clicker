package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	boxes []*countBox

	mouseEvent   tea.MouseEvent
	clickStarted bool

	w window

	debugMsg string
	ready    bool
}

type box struct {
	x      int
	y      int
	width  int
	height int
	row    int
	col    int

	text     string
	rendered string

	m *model
}

type countBox struct {
	box
	count int
}

type window struct {
	width  int
	height int
}

func (b *box) String() string {
	return b.rendered
}

func (b *box) Update(text string) {
	rendered := lipgloss.
		NewStyle().
		Align(lipgloss.Center).
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Margin(1, 0, 0, 2).
		Padding(2).
		Width(20).
		Render(text)

	b.text = text

	b.width = lipgloss.Width(rendered) - 2
	b.height = lipgloss.Height(rendered) - 1

	x := 2
	y := 1

	col := 0
	row := 0

	for _, cb := range b.m.boxes {
		if cb.x == b.x && cb.y == b.y {
			break
		}

		possibleX := x + 2 + cb.width
		possibleY := y + 1 + cb.height

		if possibleX+b.width < b.m.w.width {
			x = possibleX
			col++
		} else {
			x = 2
			y = possibleY
			col = 0
			row++
		}
	}

	b.x = x
	b.y = y
	b.col = col
	b.row = row

	b.rendered = rendered
}

func (b *countBox) Update(count int) {
	b.count = count
	b.box.Update(fmt.Sprintf("Count: %d", count))
}

func mapCountBoxes(bs []*countBox) []string {
	var strs []string
	for _, box := range bs {
		strs = append(strs, box.rendered)
	}
	return strs
}

var initialModel = model{}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if !m.ready {
			m.w.height = msg.Height
			m.w.width = msg.Width
			m.ready = true
		} else {
			m.w.height = msg.Height
			m.w.width = msg.Width
		}
		m.newCountBox(0)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "n":
			m.newCountBox(0)
		}

	case tea.MouseMsg:
		m.mouseEvent = tea.MouseEvent(msg)
		switch m.mouseEvent.Type {
		case tea.MouseLeft:
			m.clickStarted = true
		case tea.MouseRelease:
			if m.clickStarted {
				m.clickStarted = false
				for _, b := range m.boxes {
					isInBox := m.IsInBox(&b.box)
					if isInBox {
						b.Update(b.count + 1)
					}
				}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	count := len(m.boxes)

	var rows []string
	currentRow := 0
	first := 0
	for i, b := range m.boxes {
		if (currentRow < b.row) {
			renders := mapCountBoxes(m.boxes[first:i])
			row := lipgloss.JoinHorizontal(lipgloss.Bottom, renders...)
			rows = append(rows, row)
			first = i
			currentRow++
		}
	}
	if count > first {
		renders := mapCountBoxes(m.boxes[first:])
		row := lipgloss.JoinHorizontal(lipgloss.Bottom, renders...)
		rows = append(rows, row)
	}
	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m *model) newCountBox(count int) {
	newBox := countBox{}
	newBox.m = m

	newBox.Update(count)

	m.boxes = append(m.boxes, &newBox)
}

func (m *model) IsInBox(b *box) bool {
	me := m.mouseEvent


	if me.X < b.x {
		return false
	}
	if me.X > b.x+b.width-1 {
		return false
	}
	if me.Y < b.y {
		return false
	}
	if me.Y > b.y+b.height-1 {
		return false
	}

	return true
}

func main() {
	p := tea.NewProgram(initialModel, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

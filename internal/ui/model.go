package ui

import (
	"fmt"
	"strings"

	"atlas.hash/internal/hash"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true).MarginBottom(1)
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	hashKeyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Bold(true).Width(15)
	hashValStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	matchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true).Background(lipgloss.Color("22"))
)

type state int

const (
	stateInputFile state = iota
	stateViewHashes
)

type Model struct {
	state       state
	filePath    string
	results     hash.Results
	err         error
	input       textinput.Model
	compareHash string
}

func NewModel(filePath string) Model {
	ti := textinput.New()
	m := Model{
		state:    stateInputFile,
		filePath: filePath,
		input:    ti,
	}

	if filePath != "" {
		m.state = stateViewHashes
		m.input.Placeholder = "Paste hash to compare..."
		m.input.Focus()
	} else {
		m.input.Placeholder = "Enter file path..."
		m.input.Focus()
	}

	return m
}

type hashResultMsg struct {
	res hash.Results
	err error
}

func computeHashCmd(path string) tea.Cmd {
	return func() tea.Msg {
		res, err := hash.Compute(path)
		return hashResultMsg{res: res, err: err}
	}
}

func (m Model) Init() tea.Cmd {
	if m.state == stateViewHashes {
		return tea.Batch(textinput.Blink, computeHashCmd(m.filePath))
	}
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.state == stateInputFile {
				m.filePath = strings.TrimSpace(m.input.Value())
				if m.filePath == "" {
					return m, nil
				}
				m.state = stateViewHashes
				m.input.SetValue("")
				m.input.Placeholder = "Paste hash to compare..."
				return m, computeHashCmd(m.filePath)
			} else if m.state == stateViewHashes {
				m.compareHash = strings.ToLower(strings.TrimSpace(m.input.Value()))
				// Keep focused to allow changing comparison
			}
		}

	case hashResultMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateInputFile
			m.filePath = ""
			m.input.SetValue("")
			m.input.Placeholder = "Error loading file. Enter file path..."
		} else {
			m.results = msg.res
			m.err = nil
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Atlas Hash"))
	b.WriteRune('\n')

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteRune('\n')
		b.WriteRune('\n')
	}

	if m.state == stateInputFile {
		b.WriteString("File to hash:")
		b.WriteRune('\n')
		b.WriteString(m.input.View())
		b.WriteRune('\n')
		b.WriteRune('\n')
		b.WriteString("(esc to quit)")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("Target: %s", m.filePath))
	b.WriteRune('\n')
	b.WriteRune('\n')

	if m.results.MD5 == "" && m.err == nil {
		b.WriteString("Computing hashes...")
		b.WriteRune('\n')
	} else if m.err == nil {
		renderHash := func(name, val string) {
			valStyle := hashValStyle
			nameStr := name
			if m.compareHash != "" && m.compareHash == val {
				valStyle = matchStyle
				nameStr = name + " [MATCH]"
			}
			b.WriteString(hashKeyStyle.Render(nameStr))
			b.WriteString(valStyle.Render(val))
			b.WriteRune('\n')
		}

		renderHash("MD5", m.results.MD5)
		renderHash("SHA1", m.results.SHA1)
		renderHash("SHA256", m.results.SHA256)
		renderHash("SHA512", m.results.SHA512)
	}

	b.WriteRune('\n')
	b.WriteString("Compare:")
	b.WriteRune('\n')
	b.WriteString(m.input.View())
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString("(enter to highlight matches, esc to quit)")

	return b.String()
}

package interactive

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/donaldknoller/chat-cli/internal/anthropic"
	"github.com/donaldknoller/chat-cli/internal/common_llm"
	"strings"
)

type Model struct {
	session  common_llm.Session
	viewport viewport.Model
	textarea textarea.Model
	renderer *glamour.TermRenderer
	width    int
	height   int
	lines    int
}

var (
	userStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5"))
	assistantStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	errorStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("1"))
)

const maxLines = 5
const defaultLines = 2

func InitialModel() Model {
	session := anthropic.NewSession()

	ta := textarea.New()
	ta.Focus()
	ta.Prompt = "┃ "
	ta.CharLimit = -1
	ta.SetWidth(80)
	ta.SetHeight(defaultLines)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	vp := viewport.New(80, 24)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		//glamour.WithEnvironmentConfig(),
		//glamour.WithWordWrap(0),
	)
	m := Model{
		viewport: vp,
		textarea: ta,
		renderer: renderer,
		session:  session,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch()
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	m.textarea, cmd = m.textarea.Update(message)
	cmds = append(cmds, cmd)
	m.viewport, cmd = m.viewport.Update(message)
	cmds = append(cmds, cmd)
	switch msg := message.(type) {

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return m, tea.Quit
		//case tea.KeyCtrlOpenBracket:
		//	m.viewport.LineUp(1)
		//case tea.KeyCtrlCloseBracket:
		//	m.viewport.LineDown(1)
		case tea.KeyEnter:
			if m.lines > maxLines {
				break
			}
			m.lines = m.lines + 1
			m.textarea.SetHeight(m.lines)
			m.viewport.Height = m.height - m.textarea.Height()
			m.viewport.SetContent(m.Session())
		// Submit
		case tea.KeyCtrlS:
			if m.session.IsStreaming() {
				break
			}
			submittedText := m.textarea.Value()
			if strings.TrimSpace(submittedText) == "" {
				break
			}
			m.textarea.SetHeight(defaultLines)
			m.session.AddMessage(submittedText, common_llm.User)
			cmds = append(cmds, streamMessage(m.session))
			m.viewport.SetContent(m.Session())
			m.viewport.GotoBottom()
			m.textarea.Reset()
			m.textarea.Blur()
		}
	case tea.WindowSizeMsg:
		if msg.Width == 0 || msg.Height == 0 {
			break
		}
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - m.textarea.Height()
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetContent(m.Session())
		m.viewport.GotoBottom()

	case common_llm.StreamResponse:
		if msg.Err != nil {
			m.textarea.Focus()
			break
		}
		if msg.Done {
			m.textarea.Focus()
			break
		}
		m.session.AddChunk(msg.Chunk)
		m.viewport.SetContent(m.Session())
		m.viewport.GotoBottom()

	case error:
		return m, tea.Quit
	}
	return m, tea.Batch(cmds...)
}

func (m Model) Session() string {
	ss := []string{}
	for _, st := range m.session.ChatMessages() {
		var role string
		var content string
		switch st.Role {
		case common_llm.User:
			role = userStyle.Render(st.Role.String())
		case common_llm.Assistant:
			role = assistantStyle.Render(st.Content)
		default:
			role = errorStyle.Render(st.Content)
		}
		content, _ = m.renderer.Render(st.Content)
		final := fmt.Sprintf("%s : %s", role, content)
		ss = append(ss, final)
	}

	if m.session.IsStreaming() {
		add, _ := m.renderer.Render(m.session.StreamingMessage())
		ss = append(ss, add)
	}
	s := strings.Join(ss, "\n")
	return s
}

func (m Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "..."
	}
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.viewport.View(),
		m.textarea.View(),
	)
}

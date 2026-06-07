package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"yolo-team/yolo-cli/internal/ai"
	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/handler"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------------------------------------------------------------------------
// Agent integration
// ---------------------------------------------------------------------------

// chatAgent wraps the internal/ai.Agent for use in the TUI.
type chatAgent struct {
	agent      *ai.Agent
	sessionMgr *ai.SessionManager
}

func newChatAgent(workspace string) (*chatAgent, error) {
	settings, err := config.LoadSettings()
	if err != nil {
		settings = &config.Settings{Workspace: workspace}
	}
	cfg := ai.LoadConfig(settings)
	sessionMgr := ai.NewSessionManager(workspace)

	// Create handler instances for tools.
	projectH := handler.NewProjectHandler()
	executorH := handler.NewExecutorHandler()
	issueH := handler.NewIssueHandler()
	docH := handler.NewDocumentHandler(workspace)

	agent, err := ai.NewAgent(cfg, sessionMgr, projectH, executorH, issueH, issueH, docH)
	if err != nil {
		return nil, err
	}
	return &chatAgent{agent: agent, sessionMgr: sessionMgr}, nil
}

// SessionInfo holds minimal session data for the TUI.
type SessionInfo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

// sessionItem adapts SessionInfo for bubbles/list.Model.
type sessionItem struct {
	session SessionInfo
}

func (i sessionItem) Title() string       { return i.session.Title }
func (i sessionItem) Description() string { return i.session.ID }
func (i sessionItem) FilterValue() string { return i.session.Title + " " + i.session.ID }

// ---------------------------------------------------------------------------
// Styles
// ---------------------------------------------------------------------------

var (
	subtle    = lipgloss.AdaptiveColor{Light: "#9B9B9B", Dark: "#5C5C5C"}
	highlight = lipgloss.AdaptiveColor{Light: "#FF5733", Dark: "#FF8C66"}
	special   = lipgloss.AdaptiveColor{Light: "#2070B0", Dark: "#6BB5FF"}
	success   = lipgloss.AdaptiveColor{Light: "#2E8B57", Dark: "#5FD98D"}
	warning   = lipgloss.AdaptiveColor{Light: "#D4A017", Dark: "#F0C75E"}

	sessionListWidth = 28

	appStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#3E3E3E"))

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5733")).
			Bold(true).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Padding(0, 1)

	sessionListStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(subtle).
				Padding(0, 1)

	sessionListTitleStyle = lipgloss.NewStyle().
				Foreground(highlight).
				Bold(true).
				Padding(0, 1)

	messageAreaStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(subtle).
				Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(special).
			Padding(0, 1)

	focusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(highlight).
				Padding(0, 1)

	userMsgStyle = lipgloss.NewStyle().
			Foreground(success).
			Bold(true)

	assistantMsgStyle = lipgloss.NewStyle().
				Foreground(special).
				Bold(true)

	userBubbleStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(success).
			Padding(0, 1).
			MarginBottom(1).
			Width(60)

	assistantBubbleStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(special).
				Padding(0, 1).
				MarginBottom(1).
				Width(60)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(subtle).
			Background(lipgloss.Color("#1A1A2E")).
			Padding(0, 1)

	loadingStyle = lipgloss.NewStyle().
			Foreground(warning).
			Italic(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

const (
	focusSessionList = iota
	focusMessageArea
)

// ---------------------------------------------------------------------------
// TUI Model
// ---------------------------------------------------------------------------

// chatMessage represents a single message in the TUI.
type chatMessage struct {
	Role    string
	Content string
}

type model struct {
	// Components
	sessionList list.Model
	messageView viewport.Model
	inputArea   textarea.Model
	spin        spinner.Model

	// Layout
	width  int
	height int
	ready  bool

	// Focus: 0 = session list, 1 = message area
	focusPanel int

	// Data
	agent          *chatAgent
	sessions       []SessionInfo
	sessionListMap map[string]*SessionInfo
	currentSession string // sessionID
	loading        bool

	// Messages for current session
	messages []chatMessage
}

// ---------------------------------------------------------------------------
// Messages
// ---------------------------------------------------------------------------

type chatResponseMsg struct {
	content string
	err     error
}

type sessionLoadedMsg struct {
	id       string
	messages []chatMessage
	err      error
}

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

// Run starts the TUI application.
func Run(workspace string) error {
	agent, err := newChatAgent(workspace)
	if err != nil {
		return fmt.Errorf("initialize AI agent: %w", err)
	}

	ta := textarea.New()
	ta.Placeholder = "输入消息... (Shift+Enter 换行, Enter 发送)"
	ta.ShowLineNumbers = false
	ta.CharLimit = 0
	ta.KeyMap.InsertNewline.SetEnabled(false)

	spin := spinner.New()
	spin.Style = loadingStyle
	spin.Spinner = spinner.Dot

	sessions, _ := agent.sessionMgr.GetSessions()
	sessionInfos := make([]SessionInfo, 0)
	sessionMap := make(map[string]*SessionInfo)
	var currentID string
	for _, s := range sessions {
		info := SessionInfo{ID: s.SessionID, Title: s.Title, UpdatedAt: s.UpdatedAt}
		sessionInfos = append(sessionInfos, info)
		sessionMap[s.SessionID] = &info
		if currentID == "" {
			currentID = s.SessionID
		}
	}

	items := make([]list.Item, len(sessionInfos))
	for i, si := range sessionInfos {
		items[i] = sessionItem{session: si}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "会话列表"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.DisableQuitKeybindings()
	l.Styles.Title = sessionListTitleStyle

	// Load messages for current session
	var messages []chatMessage
	if currentID != "" {
		history, err := agent.sessionMgr.GetHistory(currentID, -1)
		if err == nil {
			for _, msg := range history {
				messages = append(messages, chatMessage{Role: string(msg.Role), Content: msg.Content})
			}
		}
	}

	m := &model{
		inputArea:      ta,
		spin:           spin,
		agent:          agent,
		sessionList:    l,
		sessions:       sessionInfos,
		sessionListMap: sessionMap,
		currentSession: currentID,
		messages:       messages,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// Init implements tea.Model.
func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spin.Tick,
	)
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

// Update implements tea.Model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)

	case tea.KeyMsg:
		return m.handleKey(msg)

	case chatMessageCmd:
		return m.handleSendMessage(msg.content)

	case chatResponseMsg:
		m.loading = false
		if msg.err != nil {
			m.messages = append(m.messages, chatMessage{Role: "system", Content: fmt.Sprintf("错误: %v", msg.err)})
			m.renderMessages()
			m.messageView.GotoBottom()
			return m, nil
		}
		m.messages = append(m.messages, chatMessage{Role: "assistant", Content: msg.content})
		m.renderMessages()
		m.messageView.GotoBottom()
		return m, nil

	case sessionLoadedMsg:
		if msg.err == nil && msg.id == m.currentSession {
			m.messages = msg.messages
			m.renderMessages()
			m.messageView.GotoBottom()
		}
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spin, cmd = m.spin.Update(msg)
			return m, cmd
		}
		return m, nil
	}

	// Route messages to the focused component.
	var cmd tea.Cmd
	switch m.focusPanel {
	case focusSessionList:
		m.sessionList, cmd = m.sessionList.Update(msg)
		// Check if selected item changed.
		if selected, ok := m.sessionList.SelectedItem().(sessionItem); ok {
			if selected.session.ID != m.currentSession {
				m.switchSession(selected.session.ID)
			}
		}
	case focusMessageArea:
		m.inputArea, cmd = m.inputArea.Update(msg)
	}

	return m, cmd
}

func (m *model) handleWindowSize(msg tea.WindowSizeMsg) (*model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height

	inputHeight := 5
	statusBarHeight := 1
	headerHeight := 2
	verticalMargin := 2

	availHeight := msg.Height - headerHeight - inputHeight - statusBarHeight - verticalMargin
	listH := availHeight - 2
	if listH < 10 {
		listH = 10
	}

	m.sessionList.SetWidth(sessionListWidth)
	m.sessionList.SetHeight(listH)

	msgViewWidth := msg.Width - sessionListWidth - 6
	if msgViewWidth < 40 {
		msgViewWidth = 40
	}

	if !m.ready {
		m.messageView = viewport.New(msgViewWidth, availHeight)
		m.ready = true
	}
	m.messageView.Width = msgViewWidth
	m.messageView.Height = availHeight

	m.inputArea.SetWidth(msg.Width - 4)
	m.inputArea.SetHeight(3)

	return m, nil
}

func (m *model) handleKey(msg tea.KeyMsg) (*model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyCtrlL:
		return m.handleClearSession()

	case tea.KeyCtrlN:
		return m.handleNewSession()

	case tea.KeyTab:
		m.toggleFocus()
		return m, nil

	case tea.KeyEscape:
		m.focusPanel = focusMessageArea
		m.inputArea.Focus()
		return m, nil
	}

	if m.focusPanel == focusMessageArea && msg.Type == tea.KeyEnter && !msg.Alt {
		return m.handleUserSend()
	}

	return m, nil
}

func (m *model) toggleFocus() {
	if m.focusPanel == focusSessionList {
		m.focusPanel = focusMessageArea
		m.inputArea.Focus()
	} else {
		m.focusPanel = focusSessionList
		m.inputArea.Blur()
	}
}

func (m *model) switchSession(sessionID string) {
	m.currentSession = sessionID
	m.loading = false

	history, err := m.agent.sessionMgr.GetHistory(sessionID, -1)
	if err == nil {
		m.messages = make([]chatMessage, 0)
		for _, msg := range history {
			m.messages = append(m.messages, chatMessage{Role: string(msg.Role), Content: msg.Content})
		}
	} else {
		m.messages = nil
	}
	m.renderMessages()
	m.messageView.GotoBottom()
}

func (m *model) handleNewSession() (*model, tea.Cmd) {
	sessionID, err := m.agent.sessionMgr.CreateSession()
	if err != nil {
		m.messages = append(m.messages, chatMessage{Role: "system", Content: fmt.Sprintf("创建会话失败: %v", err)})
		m.renderMessages()
		return m, nil
	}

	info := SessionInfo{ID: sessionID, Title: "New Chat", UpdatedAt: time.Now()}
	m.sessions = append(m.sessions, info)
	m.sessionListMap[sessionID] = &m.sessions[len(m.sessions)-1]

	items := make([]list.Item, len(m.sessions))
	for i, si := range m.sessions {
		items[i] = sessionItem{session: si}
	}
	cmd := m.sessionList.SetItems(items)
	m.sessionList.Select(len(items) - 1)

	m.currentSession = sessionID
	m.messages = nil
	m.inputArea.Reset()
	m.renderMessages()
	m.messageView.SetContent("")
	m.loading = false

	return m, cmd
}

func (m *model) handleClearSession() (*model, tea.Cmd) {
	m.messages = nil
	m.renderMessages()
	return m, nil
}

func (m *model) handleUserSend() (*model, tea.Cmd) {
	if m.loading {
		return m, nil
	}
	text := strings.TrimSpace(m.inputArea.Value())
	if text == "" {
		return m, nil
	}

	// Handle slash commands
	switch {
	case text == "/sessions":
		return m.handleListSessions()
	case text == "/clear":
		return m.handleClearSession()
	case strings.HasPrefix(text, "/delete "):
		return m.handleDeleteSession(strings.TrimSpace(text[8:]))
	}

	// Ensure there's a current session.
	if m.currentSession == "" {
		sessionID, err := m.agent.sessionMgr.CreateSession()
		if err != nil {
			m.messages = append(m.messages, chatMessage{Role: "system", Content: fmt.Sprintf("创建会话失败: %v", err)})
			m.renderMessages()
			return m, nil
		}
		info := SessionInfo{ID: sessionID, Title: "New Chat", UpdatedAt: time.Now()}
		m.sessions = append(m.sessions, info)
		m.sessionListMap[sessionID] = &m.sessions[len(m.sessions)-1]
		items := make([]list.Item, len(m.sessions))
		for i, si := range m.sessions {
			items[i] = sessionItem{session: si}
		}
		cmd := m.sessionList.SetItems(items)
		m.sessionList.Select(len(items) - 1)
		m.currentSession = sessionID
		m.inputArea.Reset()
		return m, cmd
	}

	return m, func() tea.Msg { return chatMessageCmd{content: text} }
}

func (m *model) handleSendMessage(text string) (*model, tea.Cmd) {
	m.messages = append(m.messages, chatMessage{Role: "user", Content: text})
	m.renderMessages()
	m.messageView.GotoBottom()
	m.inputArea.Reset()
	m.loading = true

	return m, func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		var fullResponse string
		_, _, err := m.agent.agent.ChatStream(ctx, m.currentSession, text, func(chunk string) {
			fullResponse += chunk
		})

		if err != nil {
			return chatResponseMsg{err: err}
		}
		return chatResponseMsg{content: fullResponse}
	}
}

type chatMessageCmd struct {
	content string
}

func (m *model) handleListSessions() (*model, tea.Cmd) {
	var buf strings.Builder
	buf.WriteString("会话列表:\n\n")
	if len(m.sessions) == 0 {
		buf.WriteString("  (无会话)\n")
	} else {
		for i, s := range m.sessions {
			buf.WriteString(fmt.Sprintf("  %d. %s (%s)\n", i+1, s.Title, s.ID))
		}
		buf.WriteString("\nCtrl+N 新建会话 | Tab 切换选择会话")
	}

	// Update title display
	m.messages = append(m.messages, chatMessage{Role: "system", Content: buf.String()})
	m.renderMessages()
	m.messageView.GotoBottom()

	return m, nil
}

func (m *model) handleDeleteSession(id string) (*model, tea.Cmd) {
	if err := m.agent.sessionMgr.DeleteSession(id); err != nil {
		m.messages = append(m.messages, chatMessage{Role: "system", Content: fmt.Sprintf("删除失败: %v", err)})
		m.renderMessages()
		return m, nil
	}

	// Reload sessions.
	sessions, _ := m.agent.sessionMgr.GetSessions()
	m.sessions = make([]SessionInfo, 0)
	m.sessionListMap = make(map[string]*SessionInfo)
	for _, s := range sessions {
		info := SessionInfo{ID: s.SessionID, Title: s.Title, UpdatedAt: s.UpdatedAt}
		m.sessions = append(m.sessions, info)
		m.sessionListMap[s.SessionID] = &m.sessions[len(m.sessions)-1]
	}

	items := make([]list.Item, len(m.sessions))
	for i, si := range m.sessions {
		items[i] = sessionItem{session: si}
	}
	cmd := m.sessionList.SetItems(items)

	if id == m.currentSession {
		if len(m.sessions) > 0 {
			m.currentSession = m.sessions[0].ID
			m.switchSession(m.currentSession)
		} else {
			m.currentSession = ""
			m.messages = nil
			m.renderMessages()
		}
	}

	return m, cmd
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

// View implements tea.Model.
func (m *model) View() string {
	if !m.ready {
		return "加载中..."
	}

	var buf strings.Builder

	// Header
	title := titleStyle.Render(" Yolo-Team Chat ")
	buf.WriteString(lipgloss.JoinHorizontal(lipgloss.Top,
		title,
		helpStyle.Render("Ctrl+N 新建 | Tab 切换 | Ctrl+L 清屏 | Ctrl+C 退出"),
	))
	buf.WriteString("\n")

	// Main content: session list + message view
	sessionPane := m.renderSessionPane()
	messagePane := m.renderMessagePane()
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, sessionPane, messagePane)
	buf.WriteString(mainContent)
	buf.WriteString("\n")

	// Input area
	inputView := m.inputArea.View()
	if m.focusPanel == focusMessageArea {
		inputView = focusedInputStyle.Render(inputView)
	} else {
		inputView = inputStyle.Render(inputView)
	}
	buf.WriteString(inputView)
	buf.WriteString("\n")

	// Status bar
	statusText := m.renderStatus()
	buf.WriteString(statusBarStyle.Render(statusText))

	return appStyle.Width(m.width - 2).Render(buf.String())
}

func (m *model) renderSessionPane() string {
	if len(m.sessions) == 0 {
		return sessionListStyle.
			Width(sessionListWidth).
			Height(m.messageView.Height).
			Render("  无会话\n  Ctrl+N 新建")
	}

	listView := m.sessionList.View()
	style := sessionListStyle.Width(sessionListWidth).Height(m.messageView.Height + 2)
	if m.focusPanel == focusSessionList {
		style = style.BorderForeground(highlight)
	}
	return style.Render(listView)
}

func (m *model) renderMessagePane() string {
	width := m.width - sessionListWidth - 8
	if width < 40 {
		width = 40
	}

	if m.currentSession == "" {
		content := lipgloss.NewStyle().
			Width(width).
			Height(m.messageView.Height).
			Align(lipgloss.Center).
			Foreground(subtle).
			Render("\n\n欢迎使用 Yolo-Team Chat\nCtrl+N 新建会话开始聊天")
		m.messageView.SetContent(content)
	} else if len(m.messages) == 0 {
		welcome := lipgloss.NewStyle().
			Foreground(subtle).
			Render(" 开始对话吧！输入消息后按 Enter 发送。")
		m.messageView.SetContent(welcome)
	} else {
		m.messageView.SetContent(m.messagesContent())
	}

	style := messageAreaStyle.
		Width(width + 2).
		Height(m.messageView.Height + 2)
	if m.focusPanel == focusMessageArea {
		style = style.BorderForeground(highlight)
	}

	return style.Render(m.messageView.View())
}

func (m *model) renderStatus() string {
	var parts []string
	if m.focusPanel == focusSessionList {
		parts = append(parts, "[会话列表] (Tab 切换)")
	} else {
		parts = append(parts, "[输入模式]")
	}
	if m.loading {
		parts = append(parts, m.spin.View()+" AI 思考中...")
	}
	if info, ok := m.sessionListMap[m.currentSession]; ok {
		parts = append(parts, fmt.Sprintf("会话: %s", info.Title))
	}
	return strings.Join(parts, "  |  ")
}

func (m *model) messagesContent() string {
	var buf strings.Builder
	viewWidth := m.messageView.Width
	if viewWidth < 40 {
		viewWidth = 40
	}

	for i, msg := range m.messages {
		switch msg.Role {
		case "user":
			header := userMsgStyle.Render("  You  ")
			content := lipgloss.NewStyle().Width(viewWidth - 6).Render(msg.Content)
			bubble := lipgloss.JoinVertical(lipgloss.Left, header, content)
			buf.WriteString(userBubbleStyle.Width(viewWidth).Render(bubble))
		case "assistant":
			header := assistantMsgStyle.Render("  AI  ")
			content := lipgloss.NewStyle().Width(viewWidth - 6).Render(msg.Content)
			bubble := lipgloss.JoinVertical(lipgloss.Left, header, content)
			buf.WriteString(assistantBubbleStyle.Width(viewWidth).Render(bubble))
		default:
			buf.WriteString(lipgloss.NewStyle().
				Foreground(subtle).
				Italic(true).
				Width(viewWidth).
				Render(msg.Content))
		}
		if i < len(m.messages)-1 {
			buf.WriteString("\n\n")
		}
	}
	return buf.String()
}

func (m *model) renderMessages() {
	m.messageView.SetContent(m.messagesContent())
}

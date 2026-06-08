package tui

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"yolo-team/yolo-cli/internal/ai"
	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/handler"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ---------------------------------------------------------------------------
// Agent integration
// ---------------------------------------------------------------------------

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

// ---------------------------------------------------------------------------
// TUI Model
// ---------------------------------------------------------------------------

type model struct {
	messageView viewport.Model
	inputArea   textarea.Model
	spin        spinner.Model

	width  int
	height int
	ready  bool

	agent          *chatAgent
	sessions       []SessionInfo
	sessionListMap map[string]*SessionInfo
	currentSession string
	loading        bool

	// Flat transcript (one entry = one rendered line)
	transcript      []string
	transcriptDirty bool

	// Streaming state
	streamingContent strings.Builder
	streamingEmpty   bool

	program *tea.Program
}

// ---------------------------------------------------------------------------
// Initialization
// ---------------------------------------------------------------------------

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
	ta.Focus()

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

	var transcript []string
	if currentID != "" {
		history, err := agent.sessionMgr.GetHistory(currentID, -1)
		if err == nil {
			for _, msg := range history {
				if content := msg.Content; content != "" {
					rendered := renderTranscriptMessage(string(msg.Role), content, 80)
					transcript = append(transcript, rendered...)
				}
			}
		}
	}

	m := &model{
		inputArea:      ta,
		spin:           spin,
		agent:          agent,
		sessions:       sessionInfos,
		sessionListMap: sessionMap,
		currentSession: currentID,
		transcript:     transcript,
	}

	// Redirect logs to file so stderr doesn't corrupt the TUI layout.
	logFile, err := os.OpenFile(filepath.Join(workspace, "tui.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.SetOutput(logFile)
		defer log.SetOutput(os.Stderr)
		defer logFile.Close()
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	m.program = p
	_, err = p.Run()
	return err
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		m.spin.Tick,
	)
}

// ---------------------------------------------------------------------------
// Layout helpers
// ---------------------------------------------------------------------------

func (m *model) boxWidth() int {
	w := m.width - 4
	if w < 40 {
		w = 40
	}
	return w
}

// bottomRows returns the height of all fixed bottom elements.
func (m *model) bottomRows() int {
	rows := 0
	// Working line (spinner) when loading
	if m.loading {
		rows++
	}
	// Input box: textarea height + top border + bottom border
	rows += 3 + 2 // 3 rows textarea + 2 borders
	// Status line
	rows++
	return rows
}

// transcriptHeight returns the height available for the transcript viewport.
func (m *model) transcriptHeight() int {
	// Header: title+border = 3, rest = m.height - bottomRows - margins
	h := m.height - m.bottomRows() - 4
	if h < 5 {
		h = 5
	}
	return h
}

// ---------------------------------------------------------------------------
// Transcript management
// ---------------------------------------------------------------------------

func (m *model) flushTranscript() {
	if !m.transcriptDirty {
		return
	}
	content := strings.Join(m.transcript, "\n")
	m.messageView.SetContent(content)
	m.transcriptDirty = false
}

func (m *model) appendUserMessage(text string) {
	w := m.boxWidth()
	lines := renderTranscriptMessage("user", text, w)
	m.transcript = append(m.transcript, lines...)
	m.transcriptDirty = true
}

func (m *model) beginAssistantMessage() {
	m.streamingContent.Reset()
	m.streamingEmpty = true
	m.transcript = append(m.transcript, "")
	m.transcriptDirty = true
}

func (m *model) renderStreamingLine() {
	content := m.streamingContent.String()
	if content == "" {
		if m.streamingEmpty {
			m.transcript[len(m.transcript)-1] = spinnerLine(m.spin.View())
			m.transcriptDirty = true
			m.flushTranscript()
			m.messageView.GotoBottom()
		}
		return
	}
	m.streamingEmpty = false
	w := m.boxWidth()
	rendered := mdRender(content, w)
	lines := strings.Split(rendered, "\n")
	if len(lines) == 0 {
		return
	}
	lastLine := lines[len(lines)-1]
	if lastLine == "" && len(lines) > 1 {
		lastLine = lines[len(lines)-2]
	}
	m.transcript[len(m.transcript)-1] = lastLine
	m.transcriptDirty = true
	m.flushTranscript()
	m.messageView.GotoBottom()
}

func (m *model) commitAssistantMessage() {
	content := m.streamingContent.String()
	if content == "" || m.streamingEmpty {
		if len(m.transcript) > 0 {
			m.transcript = m.transcript[:len(m.transcript)-1]
		}
		return
	}
	w := m.boxWidth()
	rendered := strings.Split(mdRender(content, w), "\n")
	if len(m.transcript) > 0 {
		m.transcript = m.transcript[:len(m.transcript)-1]
	}
	m.transcript = append(m.transcript, rendered...)
	m.transcriptDirty = true
	m.streamingContent.Reset()
	m.streamingEmpty = true
}

func (m *model) rebuildTranscript() {
	w := m.boxWidth()
	m.transcript = nil
	history, err := m.agent.sessionMgr.GetHistory(m.currentSession, -1)
	if err == nil {
		for _, msg := range history {
			if content := msg.Content; content != "" {
				rendered := renderTranscriptMessage(string(msg.Role), content, w)
				m.transcript = append(m.transcript, rendered...)
			}
		}
	}
	m.transcriptDirty = true
	m.flushTranscript()
	m.messageView.GotoBottom()
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		bw := m.boxWidth()
		th := m.transcriptHeight()

		if !m.ready {
			m.messageView = viewport.New(bw, th)
			m.ready = true
		}
		m.messageView.Width = bw
		m.messageView.Height = th

		m.inputArea.SetWidth(m.width - 8)
		m.inputArea.SetHeight(3)

		// Re-render with new width
		if m.currentSession != "" {
			m.rebuildTranscript()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyCtrlL:
			m.handleClearSession()
			return m, nil
		case tea.KeyCtrlN:
			return m.handleNewSession()
		case tea.KeyEnter:
			if !msg.Alt {
				return m.handleUserSend()
			}
		}
		var cmd tea.Cmd
		m.inputArea, cmd = m.inputArea.Update(msg)
		return m, cmd

	case chatMessageCmd:
		return m.handleSendMessage(msg.content)

	case chatTokenMsg:
		if msg.text != "" {
			m.streamingContent.WriteString(msg.text)
			m.renderStreamingLine()
		}
		return m, nil

	case chatResponseMsg:
		m.loading = false
		if msg.err != nil {
			m.transcript = append(m.transcript, dimLine("错误: "+msg.err.Error()))
			m.transcriptDirty = true
			m.flushTranscript()
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

	var cmd tea.Cmd
	m.inputArea, cmd = m.inputArea.Update(msg)
	return m, cmd
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

func (m *model) handleNewSession() (*model, tea.Cmd) {
	sessionID, err := m.agent.sessionMgr.CreateSession()
	if err != nil {
		m.transcript = append(m.transcript, dimLine("创建会话失败: "+err.Error()))
		m.transcriptDirty = true
		m.flushTranscript()
		return m, nil
	}

	info := SessionInfo{ID: sessionID, Title: "New Chat", UpdatedAt: time.Now()}
	m.sessions = append(m.sessions, info)
	m.sessionListMap[sessionID] = &m.sessions[len(m.sessions)-1]

	m.currentSession = sessionID
	m.transcript = nil
	m.streamingContent.Reset()
	m.inputArea.Reset()
	m.messageView.SetContent("")
	m.loading = false

	return m, nil
}

func (m *model) handleClearSession() {
	m.transcript = nil
	m.streamingContent.Reset()
	m.transcriptDirty = true
	m.flushTranscript()
}

func (m *model) handleUserSend() (*model, tea.Cmd) {
	if m.loading {
		return m, nil
	}
	text := strings.TrimSpace(m.inputArea.Value())
	if text == "" {
		return m, nil
	}

	switch {
	case text == "/help":
		return m.handleHelp()
	case text == "/sessions":
		return m.handleListSessions()
	case text == "/clear":
		m.handleClearSession()
		return m, nil
	case text == "/projects":
		return m.handleRunTool("list_projects", "{}")
	case text == "/executors":
		return m.handleRunTool("list_executors", "{}")
	case strings.HasPrefix(text, "/issues "):
		projectKey := strings.TrimSpace(text[8:])
		return m.handleRunTool("list_issues", fmt.Sprintf(`{"project_id":"%s"}`, projectKey))
	case strings.HasPrefix(text, "/delete "):
		return m.handleDeleteSession(strings.TrimSpace(text[8:]))
	}

	if m.currentSession == "" {
		sessionID, err := m.agent.sessionMgr.CreateSession()
		if err != nil {
			m.transcript = append(m.transcript, dimLine("创建会话失败: "+err.Error()))
			m.transcriptDirty = true
			m.flushTranscript()
			return m, nil
		}
		info := SessionInfo{ID: sessionID, Title: "New Chat", UpdatedAt: time.Now()}
		m.sessions = append(m.sessions, info)
		m.sessionListMap[sessionID] = &m.sessions[len(m.sessions)-1]
		m.currentSession = sessionID
		m.inputArea.Reset()
		return m, nil
	}

	return m, func() tea.Msg { return chatMessageCmd{content: text} }
}

func (m *model) handleSendMessage(text string) (*model, tea.Cmd) {
	m.appendUserMessage(text)
	m.beginAssistantMessage()
	m.flushTranscript()
	m.messageView.GotoBottom()
	m.inputArea.Reset()
	m.loading = true

	prog := m.program
	sessionID := m.currentSession
	agent := m.agent.agent

	return m, func() tea.Msg {
		prog.Send(chatTokenMsg{text: ""})

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		_, _, err := agent.ChatStream(ctx, sessionID, text, func(chunk string) {
			prog.Send(chatTokenMsg{text: chunk})
		})

		if err != nil {
			return chatResponseMsg{err: err}
		}
		prog.Send(chatTokenMsg{text: ""})
		return chatResponseMsg{}
	}
}

func (m *model) handleHelp() (*model, tea.Cmd) {
	var buf strings.Builder
	buf.WriteString("可用命令:\n\n")
	buf.WriteString("/help        显示此帮助\n")
	buf.WriteString("/new         新建会话 (或 Ctrl+N)\n")
	buf.WriteString("/sessions    列出会话\n")
	buf.WriteString("/clear       清屏 (或 Ctrl+L)\n")
	buf.WriteString("/delete <id> 删除会话\n")
	buf.WriteString("/projects    列出所有项目\n")
	buf.WriteString("/executors   列出所有执行人\n")
	buf.WriteString("/issues <project_key> 列出项目的 Issue\n")
	buf.WriteString("\nCtrl+C 退出")
	m.transcript = append(m.transcript, dimLine(buf.String()))
	m.transcriptDirty = true
	m.flushTranscript()
	m.messageView.GotoBottom()
	return m, nil
}

func (m *model) handleListSessions() (*model, tea.Cmd) {
	var buf strings.Builder
	if len(m.sessions) == 0 {
		buf.WriteString("(无会话)")
	} else {
		for i, s := range m.sessions {
			buf.WriteString(fmt.Sprintf("  %d. %s (%s)\n", i+1, s.Title, s.ID))
		}
		buf.WriteString("\nCtrl+N 新建会话")
	}
	m.transcript = append(m.transcript, dimLine(buf.String()))
	m.transcriptDirty = true
	m.flushTranscript()
	m.messageView.GotoBottom()
	return m, nil
}

func (m *model) handleRunTool(toolName, arguments string) (*model, tea.Cmd) {
	if m.currentSession == "" {
		return m.handleNewSession()
	}

	m.transcript = append(m.transcript, dimLine("运行工具: "+toolName))
	m.transcriptDirty = true
	m.flushTranscript()
	m.messageView.GotoBottom()

	return m, func() tea.Msg {
		resp := fmt.Sprintf("<tool_call data-tool-name=\"%s\">\n%s\n</tool_call>\n", toolName, arguments)
		resp += "<tool_result>\n"
		result, err := runToolDirectly(m.agent.agent, toolName, arguments)
		if err != nil {
			resp += fmt.Sprintf("错误: %s", err.Error())
		} else {
			resp += result
		}
		resp += "\n</tool_result>\n"
		return chatResponseMsg{content: resp}
	}
}

func runToolDirectly(*ai.Agent, string, string) (string, error) {
	return "", fmt.Errorf("direct tool execution not yet implemented")
}

func (m *model) handleDeleteSession(id string) (*model, tea.Cmd) {
	if err := m.agent.sessionMgr.DeleteSession(id); err != nil {
		m.transcript = append(m.transcript, dimLine("删除失败: "+err.Error()))
		m.transcriptDirty = true
		m.flushTranscript()
		return m, nil
	}

	sessions, _ := m.agent.sessionMgr.GetSessions()
	m.sessions = make([]SessionInfo, 0)
	m.sessionListMap = make(map[string]*SessionInfo)
	for _, s := range sessions {
		info := SessionInfo{ID: s.SessionID, Title: s.Title, UpdatedAt: s.UpdatedAt}
		m.sessions = append(m.sessions, info)
		m.sessionListMap[s.SessionID] = &m.sessions[len(m.sessions)-1]
	}

	if id == m.currentSession {
		if len(m.sessions) > 0 {
			m.currentSession = m.sessions[0].ID
			m.rebuildTranscript()
		} else {
			m.currentSession = ""
			m.transcript = nil
			m.messageView.SetContent("")
		}
	}

	return m, nil
}

// ---------------------------------------------------------------------------
// View
// ---------------------------------------------------------------------------

func (m *model) View() string {
	if !m.ready {
		return "加载中..."
	}

	var buf strings.Builder

	// Header
	buf.WriteString(titleStyle.Render(" Yolo-Team Chat "))
	buf.WriteString(helpStyle.Render("/help 命令 | Ctrl+N 新建 | Ctrl+C 退出"))
	buf.WriteString("\n")

	// Transcript area
	m.flushTranscript()
	tw := m.messageView.Width
	th := m.messageView.Height
	area := transcriptAreaStyle.Width(tw + 2).Height(th + 2).Render(m.messageView.View())
	buf.WriteString(area)
	buf.WriteString("\n")

	// Working line (spinner + time)
	if m.loading {
		working := workingStyle.Render(m.spin.View() + " AI 思考中...")
		bw := m.boxWidth()
		buf.WriteString(statusBlockStyle.Width(bw).MaxWidth(bw).Render(working))
		buf.WriteString("\n")
	}

	// Input box (reasonix style: top+bottom borders only)
	boxW := m.width - 4
	if boxW < 40 {
		boxW = 40
	}
	inputView := inputBoxStyle.Width(boxW).Render(m.inputArea.View())
	buf.WriteString(inputView)
	buf.WriteString("\n")

	// Status line
	var statusParts []string
	if info, ok := m.sessionListMap[m.currentSession]; ok {
		statusParts = append(statusParts, "会话: "+info.Title)
	}
	statusText := strings.Join(statusParts, "  ·  ")
	buf.WriteString(statusBlockStyle.Width(boxW).MaxWidth(boxW).Render(statusText))

	return appStyle.Width(m.width - 2).Render(buf.String())
}

// ---------------------------------------------------------------------------
// Formatted line helpers
// ---------------------------------------------------------------------------

var ansiBoldT = "\033[1m"
var ansiDimT = "\033[2m"
var ansiItalicT = "\033[3m"
var ansiResetT = "\033[0m"

func dimLine(text string) string {
	return ansiDimT + text + ansiResetT
}

func spinnerLine(spin string) string {
	return ansiDimT + spin + " AI 思考中..." + ansiResetT
}

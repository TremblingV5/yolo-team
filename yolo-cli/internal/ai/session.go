package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/cloudwego/eino/schema"
)

const (
	sessionsDir      = "sessions"
	sessionsFile     = "sessions.json"
	maxHistoryPerSession = 200
)

// Message represents a single chat message in session history.
type Message struct {
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// SessionMeta holds metadata for a session.
type SessionMeta struct {
	SessionID    string    `json:"session_id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
}

type sessionsIndex struct {
	Sessions map[string]*SessionMeta `json:"sessions"`
}

// SessionManager manages chat sessions with file-based persistence.
type SessionManager struct {
	workspace string
	mu        sync.RWMutex
	counter   int
}

// NewSessionManager creates a new session manager rooted at {workspace}/sessions/.
func NewSessionManager(workspace string) *SessionManager {
	dir := filepath.Join(workspace, sessionsDir)
	os.MkdirAll(dir, 0755)

	sm := &SessionManager{
		workspace: workspace,
	}

	// Load counter from existing sessions — find the highest sequence number.
	idx := sm.loadIndex()
	for _, meta := range idx.Sessions {
		var date, seq int
		if _, err := fmt.Sscanf(meta.SessionID, "YOLO-CHAT-%d-%d", &date, &seq); err == nil && seq > sm.counter {
			sm.counter = seq
		}
	}

	return sm
}

// CreateSession creates a new session and returns its session ID.
func (sm *SessionManager) CreateSession() (string, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.counter++
	now := time.Now()
	sessionID := fmt.Sprintf("YOLO-CHAT-%s-%03d", now.Format("20060102"), sm.counter)

	meta := &SessionMeta{
		SessionID:    sessionID,
		Title:        "New Chat",
		CreatedAt:    now,
		UpdatedAt:    now,
		MessageCount: 0,
	}

	idx := sm.loadIndex()
	idx.Sessions[sessionID] = meta
	if err := sm.saveIndex(idx); err != nil {
		return "", fmt.Errorf("save index: %w", err)
	}

	return sessionID, nil
}

// AddMessage appends a message to the session's JSONL file and updates index metadata.
func (sm *SessionManager) AddMessage(sessionID, role, content string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	idx := sm.loadIndex()
	meta, ok := idx.Sessions[sessionID]
	if !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	line, err := json.Marshal(Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	filePath := sm.messageFilePath(sessionID)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open message file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(append(line, '\n')); err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	meta.MessageCount++
	meta.UpdatedAt = time.Now()
	// Auto-set title from first user message.
	if meta.Title == "New Chat" && role == "user" {
		title := content
		if len(title) > 50 {
			title = title[:50] + "..."
		}
		meta.Title = title
	}

	// Enforce history limit.
	if meta.MessageCount > maxHistoryPerSession {
		if err := sm.trimHistory(filePath, meta.MessageCount-maxHistoryPerSession); err != nil {
			return fmt.Errorf("trim history: %w", err)
		}
		meta.MessageCount = maxHistoryPerSession
	}

	return sm.saveIndex(idx)
}

// GetHistory reads messages from a session's JSONL file and returns them as schema.Message.
// maxRounds limits the number of conversation rounds (2 messages per round); values < 0 return all.
func (sm *SessionManager) GetHistory(sessionID string, maxRounds int) ([]*schema.Message, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	filePath := sm.messageFilePath(sessionID)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("read messages: %w", err)
	}

	// Parse all stored messages.
	var stored []Message
	for _, line := range splitLines(string(data)) {
		if line == "" {
			continue
		}
		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		stored = append(stored, msg)
	}

	// Apply round limit.
	if maxRounds >= 0 && len(stored) > maxRounds*2 {
		stored = stored[len(stored)-maxRounds*2:]
	}

	// Convert to schema.Message.
	msgs := make([]*schema.Message, 0, len(stored))
	for i := range stored {
		msgs = append(msgs, &schema.Message{
			Role:    schema.RoleType(stored[i].Role),
			Content: stored[i].Content,
		})
	}

	return msgs, nil
}

// GetSessions returns all session metadata, ordered by updated_at descending.
func (sm *SessionManager) GetSessions() ([]SessionMeta, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	idx := sm.loadIndex()
	result := make([]SessionMeta, 0, len(idx.Sessions))
	for _, meta := range idx.Sessions {
		result = append(result, *meta)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UpdatedAt.After(result[j].UpdatedAt)
	})

	return result, nil
}

// DeleteSession removes a session's data file and its index entry.
func (sm *SessionManager) DeleteSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	idx := sm.loadIndex()
	if _, ok := idx.Sessions[sessionID]; !ok {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Remove message file.
	if err := os.Remove(sm.messageFilePath(sessionID)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove message file: %w", err)
	}
	delete(idx.Sessions, sessionID)

	return sm.saveIndex(idx)
}

// ResetAll removes all sessions and their data files.
func (sm *SessionManager) ResetAll() error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	dir := filepath.Join(sm.workspace, sessionsDir)
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	// Reset counter.
	sm.counter = 0

	return nil
}

// sessionsIndexPath returns the path to the sessions.json index file.
func (sm *SessionManager) sessionsIndexPath() string {
	return filepath.Join(sm.workspace, sessionsDir, sessionsFile)
}

// messageFilePath returns the path to a session's JSONL message file.
func (sm *SessionManager) messageFilePath(sessionID string) string {
	return filepath.Join(sm.workspace, sessionsDir, sessionID+".jsonl")
}

// loadIndex reads the sessions index from disk. Returns an empty index if the file doesn't exist.
func (sm *SessionManager) loadIndex() sessionsIndex {
	idx := sessionsIndex{Sessions: make(map[string]*SessionMeta)}
	data, err := os.ReadFile(sm.sessionsIndexPath())
	if err != nil {
		return idx
	}
	json.Unmarshal(data, &idx)
	if idx.Sessions == nil {
		idx.Sessions = make(map[string]*SessionMeta)
	}
	return idx
}

// saveIndex writes the sessions index to disk.
func (sm *SessionManager) saveIndex(idx sessionsIndex) error {
	dir := filepath.Join(sm.workspace, sessionsDir)
	os.MkdirAll(dir, 0755)

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(sm.sessionsIndexPath(), data, 0644)
}

// trimHistory removes the first n lines from a JSONL file.
func (sm *SessionManager) trimHistory(filePath string, n int) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := splitLines(string(data))
	if n >= len(lines) {
		return os.WriteFile(filePath, []byte{}, 0644)
	}

	remaining := lines[n:]
	var out []byte
	for _, line := range remaining {
		if line == "" {
			continue
		}
		out = append(out, line...)
		out = append(out, '\n')
	}

	return os.WriteFile(filePath, out, 0644)
}

// splitLines splits a string into lines.
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

package tui

import (
	"strings"
	"time"
)

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

// chatMessage represents a single message in the TUI.
type chatMessage struct {
	Role    string
	Content string
}

// msgSegment represents a parsed part of a message for rendering.
type msgSegment struct {
	kind     string // "text", "tool_call", "tool_result"
	content  string
	toolName string
}

// parseMessageSegments splits message text by <tool_call> / <tool_result> tags.
func parseMessageSegments(text string) []msgSegment {
	var segs []msgSegment
	for {
		ti := strings.Index(text, "<tool_call")
		ri := strings.Index(text, "<tool_result")
		idx := -1
		isToolCall := false
		if ti >= 0 && (ri < 0 || ti < ri) {
			idx = ti
			isToolCall = true
		} else if ri >= 0 {
			idx = ri
		}
		if idx < 0 {
			if text != "" {
				segs = append(segs, msgSegment{kind: "text", content: text})
			}
			break
		}
		if idx > 0 {
			segs = append(segs, msgSegment{kind: "text", content: text[:idx]})
		}
		text = text[idx:]

		if isToolCall {
			endTag := "</tool_call>"
			closeIdx := strings.Index(text, endTag)
			if closeIdx < 0 {
				segs = append(segs, msgSegment{kind: "text", content: text})
				break
			}
			toolName := ""
			nmStart := strings.Index(text, `data-tool-name="`)
			if nmStart >= 0 {
				nmStart += len(`data-tool-name="`)
				nmEnd := strings.Index(text[nmStart:], `"`)
				if nmEnd >= 0 {
					toolName = text[nmStart : nmStart+nmEnd]
				}
			}
			contentStart := strings.Index(text, ">") + 1
			if contentStart > 0 && contentStart < closeIdx {
				segs = append(segs, msgSegment{kind: "tool_call", content: text[contentStart:closeIdx], toolName: toolName})
			}
			text = text[closeIdx+len(endTag):]
		} else {
			endTag := "</tool_result>"
			closeIdx := strings.Index(text, endTag)
			if closeIdx < 0 {
				segs = append(segs, msgSegment{kind: "text", content: text})
				break
			}
			contentStart := strings.Index(text, ">") + 1
			if contentStart > 0 && contentStart < closeIdx {
				segs = append(segs, msgSegment{kind: "tool_result", content: text[contentStart:closeIdx]})
			}
			text = text[closeIdx+len(endTag):]
		}
	}
	return segs
}

// ---------------------------------------------------------------------------
// Tea message types
// ---------------------------------------------------------------------------

type chatResponseMsg struct {
	content string
	err     error
}

type chatTokenMsg struct {
	text string
}

type sessionLoadedMsg struct {
	id       string
	messages []chatMessage
	err      error
}

type chatMessageCmd struct {
	content string
}

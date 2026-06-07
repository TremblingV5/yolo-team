package ai

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"yolo-team/yolo-cli/internal/common"

	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	agent      *Agent
	sessionMgr *SessionManager
	workspace  string
}

func NewAIHandler(agent *Agent, sessionMgr *SessionManager, workspace string) *AIHandler {
	return &AIHandler{agent: agent, sessionMgr: sessionMgr, workspace: workspace}
}

// --- Chat (SSE, not using Wrap) ---

type ChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
}

// Chat godoc
//
//	@Id				Chat
//	@Summary		Chat with AI
//	@Description	Send a message to AI and receive streaming SSE response
//	@Tags			AI
//	@Accept			json
//	@Produce		text/event-stream
//	@Param			chat	body		ChatRequest	true	"Chat request"
//	@Success		200		{object}	common.Response
//	@Failure		400		{object}	common.Response
//	@Router			/api/v1/ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	if h.agent == nil {
		log.Printf("[AI] chat request rejected: agent not initialized")
		c.JSON(http.StatusInternalServerError, common.Response{Code: 50001, Message: "AI agent not initialized"})
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[AI] chat request bind error: %v", err)
		c.JSON(http.StatusBadRequest, common.Response{Code: 40001, Message: err.Error()})
		return
	}

	log.Printf("[AI] chat request: session_id=%q message=%q", req.SessionID, truncate(req.Message, 80))

	// SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	tokenCh := make(chan string)

	go func() {
		defer close(tokenCh)

		sessionID, response, err := h.agent.ChatStream(c.Request.Context(), req.SessionID, req.Message, func(token string) {
			event := fmt.Sprintf(`{"type":"token","content":"%s"}`, jsonEscape(token))
			select {
			case tokenCh <- event:
			case <-c.Request.Context().Done():
			}
		})

		log.Printf("[AI] chat stream done: session_id=%q response_len=%d err=%v", sessionID, len(response), err)

		if err != nil {
			log.Printf("[AI] chat stream error: %v", err)
			select {
			case tokenCh <- fmt.Sprintf(`{"type":"error","content":"%s"}`, jsonEscape(err.Error())):
			case <-c.Request.Context().Done():
			}
		}

		doneEvent := fmt.Sprintf(`{"type":"done","session_id":"%s"}`, sessionID)
		select {
		case tokenCh <- doneEvent:
		case <-c.Request.Context().Done():
		}
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case token, ok := <-tokenCh:
			if !ok {
				return false
			}
			fmt.Fprintf(w, "data: %s\n\n", token)
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}

// --- ListSessions ---

type ListSessionsReq struct{}

type SessionResponse struct {
	SessionID    string    `json:"session_id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	MessageCount int       `json:"message_count"`
}

// ListSessions godoc
//
//	@Id				ListSessions
//	@Summary		List AI chat sessions
//	@Description	Get all chat sessions
//	@Tags			AI
//	@Produce		json
//	@Success		200	{object}	common.Response{data=[]SessionResponse}
//	@Failure		500	{object}	common.Response
//	@Router			/api/v1/ai/sessions [get]
func (h *AIHandler) ListSessions(ctx context.Context, req ListSessionsReq) ([]SessionResponse, error) {
	sessions, err := h.sessionMgr.GetSessions()
	if err != nil {
		return nil, common.NewAppError(50001, err.Error())
	}
	result := make([]SessionResponse, len(sessions))
	for i, s := range sessions {
		result[i] = SessionResponse{
			SessionID:    s.SessionID,
			Title:        s.Title,
			CreatedAt:    s.CreatedAt,
			UpdatedAt:    s.UpdatedAt,
			MessageCount: s.MessageCount,
		}
	}
	return result, nil
}

// --- GetSessionMessages ---

type GetSessionMessagesReq struct {
	SessionID string `uri:"session_id"`
}

// GetSessionMessages godoc
//
//	@Id				GetSessionMessages
//	@Summary		Get session messages
//	@Description	Get all messages of a chat session
//	@Tags			AI
//	@Produce		json
//	@Param			session_id	path	string	true	"Session ID"
//	@Success		200	{object}	common.Response{data=[]schema.Message}
//	@Failure		404	{object}	common.Response
//	@Router			/api/v1/ai/sessions/{session_id} [get]
func (h *AIHandler) GetSessionMessages(ctx context.Context, req GetSessionMessagesReq) ([]*schema.Message, error) {
	messages, err := h.sessionMgr.GetHistory(req.SessionID, -1)
	if err != nil {
		return nil, common.NewAppError(40401, "session not found")
	}
	return messages, nil
}

// --- DeleteSession ---

type DeleteSessionReq struct {
	SessionID string `uri:"session_id"`
}

// DeleteSession godoc
//
//	@Id				DeleteSession
//	@Summary		Delete a chat session
//	@Description	Delete a session by its ID
//	@Tags			AI
//	@Produce		json
//	@Param			session_id	path		string	true	"Session ID"
//	@Success		200			{object}	common.Response
//	@Failure		404			{object}	common.Response
//	@Router			/api/v1/ai/sessions/{session_id} [delete]
func (h *AIHandler) DeleteSession(ctx context.Context, req DeleteSessionReq) (struct{}, error) {
	if err := h.sessionMgr.DeleteSession(req.SessionID); err != nil {
		return struct{}{}, common.NewAppError(40401, "session not found")
	}
	return struct{}{}, nil
}

// --- ResetSessions ---

type ResetSessionsReq struct{}

// ResetSessions godoc
//
//	@Id				ResetSessions
//	@Summary		Reset all chat sessions
//	@Description	Delete all chat sessions
//	@Tags			AI
//	@Produce		json
//	@Success		200	{object}	common.Response
//	@Failure		500	{object}	common.Response
//	@Router			/api/v1/ai/sessions/reset [post]
func (h *AIHandler) ResetSessions(ctx context.Context, req ResetSessionsReq) (struct{}, error) {
	if err := h.sessionMgr.ResetAll(); err != nil {
		return struct{}{}, common.NewAppError(50001, err.Error())
	}
	return struct{}{}, nil
}

func jsonEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

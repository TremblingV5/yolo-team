package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`{"error":"marshal failed: %s"}`, err.Error())
	}
	return string(data)
}

// ---------- Session Tools ----------

type sessionListTool struct{ sm *SessionManager }

func (t *sessionListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_sessions", Desc: "列出所有会话（按更新时间倒序）"}, nil
}
func (t *sessionListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	sessions, err := t.sm.GetSessions()
	if err != nil {
		return fmt.Sprintf("错误：获取会话列表失败 - %s", err.Error()), nil
	}
	return toJSON(sessions), nil
}

type sessionGetTool struct{ sm *SessionManager }

func (t *sessionGetTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "get_session", Desc: "获取某个会话的消息历史。必需参数: session_id（会话ID）。可选: limit（消息条数限制）"}, nil
}
func (t *sessionGetTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		SessionID string `json:"session_id"`
		Limit     string `json:"limit"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	messages, err := t.sm.GetHistory(params.SessionID, -1)
	if err != nil {
		return "", err
	}
	return toJSON(messages), nil
}

type sessionDeleteTool struct{ sm *SessionManager }

func (t *sessionDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_session", Desc: "删除一个会话及其消息历史。必需参数: session_id（会话ID）"}, nil
}
func (t *sessionDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if err := t.sm.DeleteSession(params.SessionID); err != nil {
		return fmt.Sprintf("错误：删除会话失败 - %s", err.Error()), nil
	}
	return `{"deleted":true}`, nil
}

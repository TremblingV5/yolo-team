package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/handler"
	"yolo-team/yolo-cli/internal/model"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ---------- Issue Tools ----------

type issueListTool struct{ h *handler.IssueHandler }

func (t *issueListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_issues", Desc: "列出所有 Issue。可选: project_id（项目数字ID）, status（状态: created/in_progress/done/archived）, executor_name（执行人名称）"}, nil
}
func (t *issueListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		ProjectID    json.Number `json:"project_id"`
		Status       string      `json:"status"`
		ExecutorName string      `json:"executor_name"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	req := handler.ListIssueReq{}
	if params.ProjectID.String() != "" {
		if id, err := params.ProjectID.Int64(); err == nil {
			req.ProjectID = &id
		}
	}
	if params.Status != "" {
		req.Status = &params.Status
	}
	if params.ExecutorName != "" {
		req.ExecutorName = params.ExecutorName
	}
	issues, err := t.h.List(ctx, req)
	if err != nil {
		return fmt.Sprintf("错误：获取 Issue 列表失败 - %s", err.Error()), nil
	}
	return toJSON(issues), nil
}

type issueCreateTool struct{ h *handler.IssueHandler }

func (t *issueCreateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "create_issue", Desc: "创建一个新的 Issue。必需参数: project_id（项目数字ID）, title（标题）。可选: description, priority（低/中/高/紧急）, executor_name（执行人名称）, deadline（截止日期，格式 YYYY-MM-DD）, parent_key（父 Issue 的 Key）"}, nil
}
func (t *issueCreateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		ProjectID    json.Number `json:"project_id"`
		Title        string      `json:"title"`
		Description  string      `json:"description"`
		Priority     string      `json:"priority"`
		ExecutorName string      `json:"executor_name"`
		Deadline     string      `json:"deadline"`
		ParentKey    string      `json:"parent_key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if params.ProjectID.String() == "" {
		return "错误：project_id 是必需参数，请提供项目ID", nil
	}
	if params.Title == "" {
		return "错误：title 是必需参数，请提供 Issue 标题", nil
	}
	projectID, err := params.ProjectID.Int64()
	if err != nil {
		return fmt.Sprintf("错误：无效的 project_id %q", params.ProjectID.String()), nil
	}
	result, err := t.h.Create(ctx, handler.CreateIssueReq{
		ProjectID:    projectID,
		Title:        params.Title,
		Description:  params.Description,
		Priority:     params.Priority,
		ExecutorName: params.ExecutorName,
		Deadline:     params.Deadline,
		ParentKey:    params.ParentKey,
	})
	if err != nil {
		return fmt.Sprintf("错误：创建 Issue 失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type issueGetTool struct{ h *handler.IssueHandler }

func (t *issueGetTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "get_issue", Desc: "根据 Key 获取 Issue 详情。必需参数: key（Issue Key，格式如 YOLO-ISSUE-{id}）"}, nil
}
func (t *issueGetTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	issue, err := t.h.Get(ctx, handler.GetIssueReq{Key: params.Key})
	if err != nil {
		return fmt.Sprintf("错误：未找到 Issue %q", params.Key), nil
	}
	return toJSON(issue), nil
}

type issueUpdateTool struct{ h *handler.IssueHandler }

func (t *issueUpdateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "update_issue", Desc: "更新 Issue 的字段。必需参数: key（Issue Key）。可选: title, description, status(created/in_progress/done/archived), priority(低/中/高/紧急), executor_name, deadline"}, nil
}
func (t *issueUpdateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key          string  `json:"key"`
		Title        *string `json:"title"`
		Description  *string `json:"description"`
		Status       *string `json:"status"`
		Priority     *string `json:"priority"`
		ExecutorName *string `json:"executor_name"`
		Deadline     *string `json:"deadline"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	req := handler.UpdateIssueReq{Key: params.Key, UpdateIssueRequest: model.UpdateIssueRequest{
		Title:        params.Title,
		Description:  params.Description,
		Status:       params.Status,
		Priority:     params.Priority,
		ExecutorName: params.ExecutorName,
	}}
	// Deadline requires *time.Time; not handled here.
	result, err := t.h.Update(ctx, req)
	if err != nil {
		return fmt.Sprintf("错误：更新 Issue 失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type issueDeleteTool struct{ h *handler.IssueHandler }

func (t *issueDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_issue", Desc: "删除一个 Issue。必需参数: key（Issue Key）"}, nil
}
func (t *issueDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if _, err := t.h.Delete(ctx, handler.DeleteIssueReq{Key: params.Key}); err != nil {
		return fmt.Sprintf("错误：未找到 Issue %q", params.Key), nil
	}
	return `{"deleted":true}`, nil
}

type issueTodoTool struct{ h *handler.IssueHandler }

func (t *issueTodoTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "todo_issues", Desc: "获取某个执行人的待办 Issue 列表，必需参数: executor_name（执行人名称）"}, nil
}
func (t *issueTodoTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		ExecutorName string `json:"executor_name"`
		Limit        string `json:"limit"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if params.ExecutorName == "" {
		return "错误：executor_name 是必需参数，请提供执行人名称", nil
	}
	limit := 20
	if params.Limit != "" {
		if l, err := parseInt(params.Limit); err == nil && l > 0 {
			limit = l
		}
	}
	issues, err := t.h.Todo(ctx, handler.TodoIssueReq{ExecutorName: params.ExecutorName, Limit: limit})
	if err != nil {
		return fmt.Sprintf("错误：获取待办列表失败 - %s", err.Error()), nil
	}
	return toJSON(issues), nil
}

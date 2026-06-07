package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ---------- Task Tools ----------

type taskListTool struct{ h *handler.IssueHandler }

func (t *taskListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_tasks", Desc: "列出某个 Issue 下的所有子任务。必需参数: issue_key（Issue 的 Key）"}, nil
}
func (t *taskListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		IssueKey string `json:"issue_key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	tasks, err := t.h.ListTasks(ctx, handler.ListTaskReq{IssueKey: params.IssueKey})
	if err != nil {
		return fmt.Sprintf("错误：未找到 Issue %q", params.IssueKey), nil
	}
	return toJSON(tasks), nil
}

type taskCreateTool struct{ h *handler.IssueHandler }

func (t *taskCreateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "create_task", Desc: "在 Issue 下创建子任务。必需参数: issue_key（Issue 的 Key）, title（子任务标题）。可选: description"}, nil
}
func (t *taskCreateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		IssueKey    string `json:"issue_key"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if params.IssueKey == "" {
		return "错误：issue_key 是必需参数，请提供任务卡片的 Key", nil
	}
	if params.Title == "" {
		return "错误：title 是必需参数，请提供子任务标题", nil
	}
	result, err := t.h.CreateTask(ctx, handler.CreateTaskReq{
		IssueKey:    params.IssueKey,
		Title:       params.Title,
		Description: params.Description,
	})
	if err != nil {
		return fmt.Sprintf("错误：创建子任务失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type taskUpdateTool struct{ h *handler.IssueHandler }

func (t *taskUpdateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "update_task", Desc: "更新子任务的字段或状态。必需参数: task_key（子任务 Key）。可选: title, description, status(created/in_progress/done)"}, nil
}
func (t *taskUpdateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		TaskKey     string  `json:"task_key"`
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Status      *string `json:"status"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	result, err := t.h.UpdateTask(ctx, handler.UpdateTaskReq{
		TaskKey:     params.TaskKey,
		Title:       params.Title,
		Description: params.Description,
		Status:      params.Status,
	})
	if err != nil {
		return fmt.Sprintf("错误：更新子任务失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type taskDeleteTool struct{ h *handler.IssueHandler }

func (t *taskDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_task", Desc: "删除一个子任务。必需参数: task_key（子任务 Key）"}, nil
}
func (t *taskDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		TaskKey string `json:"task_key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if _, err := t.h.DeleteTask(ctx, handler.DeleteTaskReq{TaskKey: params.TaskKey}); err != nil {
		return fmt.Sprintf("错误：未找到子任务 %q", params.TaskKey), nil
	}
	return `{"deleted":true}`, nil
}

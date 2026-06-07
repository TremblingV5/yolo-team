package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ---------- Executor Tools ----------

type executorListTool struct{ h *handler.ExecutorHandler }

func (t *executorListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_executors", Desc: "列出所有执行人（AI员工）"}, nil
}
func (t *executorListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	executors, err := t.h.List(ctx, struct{}{})
	if err != nil {
		return fmt.Sprintf("错误：获取执行人列表失败 - %s", err.Error()), nil
	}
	return toJSON(executors), nil
}

type executorCreateTool struct{ h *handler.ExecutorHandler }

func (t *executorCreateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "create_executor", Desc: "创建新的执行人（AI员工）。必需参数: name（名称）。可选: soul（灵魂设定）"}, nil
}
func (t *executorCreateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Name string `json:"name"`
		Soul string `json:"soul"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	result, err := t.h.Create(ctx, handler.CreateExecutorReq{Name: params.Name, Soul: params.Soul})
	if err != nil {
		return fmt.Sprintf("错误：创建执行人失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type executorUpdateTool struct{ h *handler.ExecutorHandler }

func (t *executorUpdateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "update_executor", Desc: "更新执行人的角色和/或灵魂设定。必需参数: name（执行人名称）。可选: role（角色: leader/architect/developer/qa）, soul（灵魂设定）"}, nil
}
func (t *executorUpdateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Name string `json:"name"`
		Role string `json:"role"`
		Soul string `json:"soul"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	result, err := t.h.Update(ctx, handler.UpdateExecutorReq{Name: params.Name, Role: params.Role, Soul: params.Soul})
	if err != nil {
		return fmt.Sprintf("错误：更新执行人失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type executorDeleteTool struct{ h *handler.ExecutorHandler }

func (t *executorDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_executor", Desc: "删除一个执行人。必需参数: name（执行人名称）"}, nil
}
func (t *executorDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if _, err := t.h.Delete(ctx, handler.DeleteExecutorReq{Name: params.Name}); err != nil {
		return fmt.Sprintf("错误：未找到执行人 %q", params.Name), nil
	}
	return `{"deleted":true}`, nil
}

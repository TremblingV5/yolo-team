package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ---------- Project Tools ----------

type projectListTool struct{ h *handler.ProjectHandler }

func (t *projectListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_projects", Desc: "列出所有项目"}, nil
}
func (t *projectListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	projects, err := t.h.List(ctx, struct{}{})
	if err != nil {
		return fmt.Sprintf("错误：获取项目列表失败 - %s", err.Error()), nil
	}
	return toJSON(projects), nil
}

type projectCreateTool struct{ h *handler.ProjectHandler }

func (t *projectCreateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "create_project", Desc: "创建一个新项目。必需参数: name（项目名称，最长32字符）。可选: description（描述，最长128字符）"}, nil
}
func (t *projectCreateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if params.Name == "" {
		return "错误：name 是必需参数，请提供项目名称", nil
	}
	result, err := t.h.Create(ctx, handler.CreateProjectReq{Name: params.Name, Description: params.Description})
	if err != nil {
		return fmt.Sprintf("错误：创建项目失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type projectGetTool struct{ h *handler.ProjectHandler }

func (t *projectGetTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "get_project", Desc: "根据 Key 获取项目详情。必需参数: key（项目 Key，格式如 YOLO-PROJECT-{id}）"}, nil
}
func (t *projectGetTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	project, err := t.h.Get(ctx, handler.GetProjectReq{Key: params.Key})
	if err != nil {
		return fmt.Sprintf("错误：未找到项目 %q", params.Key), nil
	}
	return toJSON(project), nil
}

type projectUpdateTool struct{ h *handler.ProjectHandler }

func (t *projectUpdateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "update_project", Desc: "更新项目信息。必需参数: key（项目 Key）。可选: name（新名称）, description（新描述）"}, nil
}
func (t *projectUpdateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key         string  `json:"key"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	result, err := t.h.Update(ctx, handler.UpdateProjectReq{
		Key:         params.Key,
		Name:        params.Name,
		Description: params.Description,
	})
	if err != nil {
		return fmt.Sprintf("错误：更新项目失败 - %s", err.Error()), nil
	}
	return toJSON(result), nil
}

type projectDeleteTool struct{ h *handler.ProjectHandler }

func (t *projectDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_project", Desc: "删除一个项目。必需参数: key（项目 Key）"}, nil
}
func (t *projectDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if _, err := t.h.Delete(ctx, handler.DeleteProjectReq{Key: params.Key}); err != nil {
		return fmt.Sprintf("错误：未找到项目 %q", params.Key), nil
	}
	return `{"deleted":true}`, nil
}

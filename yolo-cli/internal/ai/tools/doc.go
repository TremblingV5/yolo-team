package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

// ---------- Document Tools ----------

type documentListTool struct{ h *handler.DocumentHandler }

func (t *documentListTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "list_documents", Desc: "列出某个项目下的所有文档。必需参数: project_key（项目 Key）"}, nil
}
func (t *documentListTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		ProjectKey string `json:"project_key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	docs, err := t.h.List(ctx, handler.ListDocReq{ProjectKey: params.ProjectKey})
	if err != nil {
		return fmt.Sprintf("错误：未找到项目 %q", params.ProjectKey), nil
	}
	return toJSON(docs), nil
}

type documentCreateTool struct{ h *handler.DocumentHandler }

func (t *documentCreateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "create_document", Desc: "在项目中创建一个新文档。必需参数: project_key（项目 Key）, title（标题）。可选: content（内容）, creator（创建者）"}, nil
}
func (t *documentCreateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		ProjectKey string `json:"project_key"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		Creator    string `json:"creator"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	resp, err := t.h.Create(ctx, handler.CreateDocReq{
		ProjectKey: params.ProjectKey,
		Title:      params.Title,
		Content:    params.Content,
		Creator:    params.Creator,
	})
	if err != nil {
		return fmt.Sprintf("错误：创建文档失败 - %s", err.Error()), nil
	}
	return toJSON(resp), nil
}

type documentGetTool struct{ h *handler.DocumentHandler }

func (t *documentGetTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "get_document", Desc: "根据 Key 获取文档详情（含正文）。必需参数: key（文档 Key，格式如 YOLO-DOC-{id}）"}, nil
}
func (t *documentGetTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	resp, err := t.h.Get(ctx, handler.GetDocReq{Key: params.Key})
	if err != nil {
		return fmt.Sprintf("错误：未找到文档 %q", params.Key), nil
	}
	return toJSON(resp), nil
}

type documentUpdateTool struct{ h *handler.DocumentHandler }

func (t *documentUpdateTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "update_document", Desc: "更新文档的标题和/或内容。必需参数: key（文档 Key）。可选: title（新标题）, content（新内容）"}, nil
}
func (t *documentUpdateTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key     string  `json:"key"`
		Title   *string `json:"title"`
		Content *string `json:"content"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	resp, err := t.h.Update(ctx, handler.UpdateDocReq{Key: params.Key, Title: params.Title, Content: params.Content})
	if err != nil {
		return fmt.Sprintf("错误：更新文档失败 - %s", err.Error()), nil
	}
	return toJSON(resp), nil
}

type documentDeleteTool struct{ h *handler.DocumentHandler }

func (t *documentDeleteTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{Name: "delete_document", Desc: "删除一个文档（含物理文件）。必需参数: key（文档 Key）"}, nil
}
func (t *documentDeleteTool) InvokableRun(ctx context.Context, arguments string, opts ...tool.Option) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(arguments), &params); err != nil {
		return "", err
	}
	if _, err := t.h.Delete(ctx, handler.DeleteDocReq{Key: params.Key}); err != nil {
		return fmt.Sprintf("错误：未找到文档 %q", params.Key), nil
	}
	return "删除成功", nil
}

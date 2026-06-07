package tools

import (
	"yolo-team/yolo-cli/internal/handler"

	"github.com/cloudwego/eino/components/tool"
)

// ToolRegistry holds all tool instances and their dependencies.
type ToolRegistry struct {
	ph  *handler.ProjectHandler
	eh  *handler.ExecutorHandler
	ih  *handler.IssueHandler
	th  *handler.IssueHandler
	dh  *handler.DocumentHandler
}

// NewToolRegistry creates a new ToolRegistry with the given handler implementations.
func NewToolRegistry(ph *handler.ProjectHandler, eh *handler.ExecutorHandler, ih *handler.IssueHandler, th *handler.IssueHandler, dh *handler.DocumentHandler) *ToolRegistry {
	return &ToolRegistry{
		ph: ph,
		eh: eh,
		ih: ih,
		th: th,
		dh: dh,
	}
}

// AllTools returns all tools as Eino InvokableTool slice.
func (r *ToolRegistry) AllTools() []tool.InvokableTool {
	return []tool.InvokableTool{
		&projectListTool{h: r.ph},
		&projectCreateTool{h: r.ph},
		&projectGetTool{h: r.ph},
		&projectUpdateTool{h: r.ph},
		&projectDeleteTool{h: r.ph},

		&issueListTool{h: r.ih},
		&issueCreateTool{h: r.ih},
		&issueGetTool{h: r.ih},
		&issueUpdateTool{h: r.ih},
		&issueDeleteTool{h: r.ih},
		&issueTodoTool{h: r.ih},

		&taskListTool{h: r.th},
		&taskCreateTool{h: r.th},
		&taskUpdateTool{h: r.th},
		&taskDeleteTool{h: r.th},

		&executorListTool{h: r.eh},
		&executorCreateTool{h: r.eh},
		&executorUpdateTool{h: r.eh},
		&executorDeleteTool{h: r.eh},

		&documentListTool{h: r.dh},
		&documentCreateTool{h: r.dh},
		&documentGetTool{h: r.dh},
		&documentUpdateTool{h: r.dh},
		&documentDeleteTool{h: r.dh},
	}
}

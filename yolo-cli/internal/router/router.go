package router

import (
	"io/fs"
	"log"
	"net/http"
	"strings"

	"yolo-team/yolo-cli/internal/ai"
	"yolo-team/yolo-cli/internal/config"
	"yolo-team/yolo-cli/internal/frontend"
	"yolo-team/yolo-cli/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func createAgent(workspace string, projectH *handler.ProjectHandler, executorH *handler.ExecutorHandler, issueH, taskH *handler.IssueHandler, docH *handler.DocumentHandler) (*ai.Agent, *ai.SessionManager) {
	settings, err := config.LoadSettings()
	if err != nil {
		settings = &config.Settings{Workspace: workspace}
	}
	cfg := ai.LoadConfig(settings)
	cfg.Workspace = workspace
	sessionMgr := ai.NewSessionManager(workspace)

	agent, agentErr := ai.NewAgent(cfg, sessionMgr, projectH, executorH, issueH, taskH, docH)
	if agentErr != nil {
		log.Printf("[AI] failed to create agent: %v", agentErr)
	} else {
		log.Printf("[AI] agent created successfully (model=%s, base_url=%s)", cfg.Model, cfg.BaseURL)
	}
	return agent, sessionMgr
}

func Setup(workspace string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		projectH := handler.NewProjectHandler()
		executorH := handler.NewExecutorHandler()
		issueH := handler.NewIssueHandler()
		taskH := issueH
		docH := handler.NewDocumentHandler(workspace)

		projects := api.Group("/projects")
		{
			projects.GET("", handler.Wrap(projectH.List))
			projects.POST("", handler.Wrap(projectH.Create))
			projects.GET("/:key", handler.Wrap(projectH.Get))
			projects.PUT("/:key", handler.Wrap(projectH.Update))
			projects.DELETE("/:key", handler.Wrap(projectH.Delete))
		}

		executors := api.Group("/executors")
		{
			executors.GET("", handler.Wrap(executorH.List))
			executors.POST("", handler.Wrap(executorH.Create))
			executors.PUT("/:name", handler.Wrap(executorH.Update))
			executors.POST("/:name/delete", handler.Wrap(executorH.Delete))
		}

		issues := api.Group("/issues")
		{
			issues.GET("", handler.Wrap(issueH.List))
			issues.POST("", handler.Wrap(issueH.Create))
			issues.GET("/todo", handler.Wrap(issueH.Todo))
			issues.GET("/:key", handler.Wrap(issueH.Get))
			issues.PUT("/:key", handler.Wrap(issueH.Update))
			issues.DELETE("/:key", handler.Wrap(issueH.Delete))
			issues.POST("/:key/documents", handler.Wrap(issueH.LinkDocument))
			issues.DELETE("/:key/documents", handler.Wrap(issueH.UnlinkDocument))
		}

		projects.GET("/:key/documents", handler.Wrap(docH.List))
		projects.POST("/:key/documents", handler.Wrap(docH.Create))
		documents := api.Group("/documents")
		{
			documents.GET("/:key", handler.Wrap(docH.Get))
			documents.PUT("/:key", handler.Wrap(docH.Update))
			documents.POST("/:key/delete", handler.Wrap(docH.Delete))
		}

		// Task routes under issues
		issues.POST("/:key/tasks", handler.Wrap(issueH.CreateTask))
		issues.GET("/:key/tasks", handler.Wrap(issueH.ListTasks))
		issues.PUT("/:key/tasks/:task_key", handler.Wrap(issueH.UpdateTask))
		issues.POST("/:key/tasks/:task_key/delete", handler.Wrap(issueH.DeleteTask))

		agent, sessionMgr := createAgent(workspace, projectH, executorH, issueH, taskH, docH)
		aiH := ai.NewAIHandler(agent, sessionMgr, workspace)
		aiGroup := api.Group("/ai")
		{
			aiGroup.POST("/chat", aiH.Chat)
			aiGroup.GET("/sessions", handler.Wrap(aiH.ListSessions))
			aiGroup.GET("/sessions/:session_id", handler.Wrap(aiH.GetSessionMessages))
			aiGroup.DELETE("/sessions/:session_id", handler.Wrap(aiH.DeleteSession))
			aiGroup.POST("/sessions/reset", handler.Wrap(aiH.ResetSessions))
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	staticFS, err := fs.Sub(frontend.FS, "client-dist")
	if err == nil {
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path

			if strings.HasPrefix(path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"code": 40401, "message": "not found"})
				return
			}

			filePath := strings.TrimPrefix(path, "/")
			if filePath == "" {
				filePath = "index.html"
			}

			data, err := fs.ReadFile(staticFS, filePath)
			if err != nil {
				data, _ = fs.ReadFile(staticFS, "index.html")
			}

			contentType := "text/html; charset=utf-8"
			if strings.HasSuffix(filePath, ".js") {
				contentType = "application/javascript; charset=utf-8"
			} else if strings.HasSuffix(filePath, ".css") {
				contentType = "text/css; charset=utf-8"
			} else if strings.HasSuffix(filePath, ".svg") {
				contentType = "image/svg+xml"
			}

			c.Data(http.StatusOK, contentType, data)
		})
	}

	return r
}

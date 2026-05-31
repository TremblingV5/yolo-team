package router

import (
	"yolo-team/yolo-cli/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(workspace string) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		projectH := handler.NewProjectHandler()
		projects := api.Group("/projects")
		{
			projects.GET("", handler.Wrap(projectH.List))
			projects.POST("", handler.Wrap(projectH.Create))
			projects.GET("/:key", handler.Wrap(projectH.Get))
			projects.PUT("/:key", handler.Wrap(projectH.Update))
			projects.DELETE("/:key", handler.Wrap(projectH.Delete))
		}

		executorH := handler.NewExecutorHandler()
		executors := api.Group("/executors")
		{
			executors.GET("", handler.Wrap(executorH.List))
			executors.POST("", handler.Wrap(executorH.Create))
			executors.DELETE("/:name", handler.Wrap(executorH.Delete))
		}

		issueH := handler.NewIssueHandler()
		issues := api.Group("/issues")
		{
			issues.GET("", handler.Wrap(issueH.List))
			issues.POST("", handler.Wrap(issueH.Create))
			issues.GET("/todo", handler.Wrap(issueH.Todo))
			issues.GET("/:key", handler.Wrap(issueH.Get))
			issues.PUT("/:key", handler.Wrap(issueH.Update))
			issues.DELETE("/:key", handler.Wrap(issueH.Delete))
		}

		docH := handler.NewDocumentHandler(workspace)
		projects.GET("/:key/documents", handler.Wrap(docH.List))
		projects.POST("/:key/documents", handler.Wrap(docH.Create))
		documents := api.Group("/documents")
		{
			documents.GET("/:key", handler.Wrap(docH.Get))
			documents.PUT("/:key", handler.Wrap(docH.Update))
			documents.DELETE("/:key", handler.Wrap(docH.Delete))
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}

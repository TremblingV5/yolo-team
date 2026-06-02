package router

import (
	"io/fs"
	"net/http"
	"strings"

	"yolo-team/yolo-cli/internal/frontend"
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
			executors.PUT("/:name", handler.Wrap(executorH.Update))
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
			issues.POST("/:key/documents", handler.Wrap(issueH.LinkDocument))
			issues.DELETE("/:key/documents", handler.Wrap(issueH.UnlinkDocument))
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

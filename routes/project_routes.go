package routes

import (
	"github.com/aruncs31s/esdcprojectmodule/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	publicProjectRoutes := r.Group("/api/projects")
	{
		publicProjectRoutes.GET("", projectHandler.GetPublicProjects)
		publicProjectRoutes.GET("/:id", projectHandler.GetProject)

	}
}

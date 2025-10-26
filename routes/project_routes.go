package routes

import (
	"github.com/aruncs31s/esdcprojectmodule/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	publicProjectRoutes := r.Group("/api/public/projects")
	{
		publicProjectRoutes.GET("", projectHandler.GetAllProjects)
		publicProjectRoutes.GET("/:id", projectHandler.GetProject)

	}
}
func RegisterPrivateProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	privateProjectRoutes := r.Group("/api/projects")
	{
		privateProjectRoutes.POST("", projectHandler.CreateProject)
		privateProjectRoutes.POST("/:id/toggle-like", projectHandler.ToggleLikeProject)
		privateProjectRoutes.GET("/:id", projectHandler.GetProject)
		privateProjectRoutes.GET("", projectHandler.GetAllProjects)
		// privateProjectRoutes.PUT("/:id", projectHandler.UpdateProject)
		// privateProjectRoutes.DELETE("/:id", projectHandler.DeleteProject)
	}
}

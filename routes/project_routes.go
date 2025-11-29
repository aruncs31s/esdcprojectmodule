package routes

import (
	"github.com/aruncs31s/esdcprojectmodule/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	publicProjectRoutes := r.Group("/api/public/projects")
	{
		publicProjectRoutes.GET("", projectHandler.GetAllProjects)
		publicProjectRoutes.GET("/filter", projectHandler.FilterProjects)
		publicProjectRoutes.GET("/:id", projectHandler.GetProject)
		publicProjectRoutes.GET("/:id/stats", projectHandler.GetProjectStats)
		publicProjectRoutes.GET("/:id/comments", projectHandler.GetComments)
		publicProjectRoutes.GET("/:id/reviews", projectHandler.GetReviews)
	}
}
func RegisterPrivateProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	privateProjectRoutes := r.Group("/api/projects")
	{
		privateProjectRoutes.POST("", projectHandler.CreateProject)
		privateProjectRoutes.PUT("/:id", projectHandler.UpdateProject)
		privateProjectRoutes.DELETE("/:id", projectHandler.DeleteProject)
		privateProjectRoutes.POST("/:id/toggle-like", projectHandler.ToggleLikeProject)
		privateProjectRoutes.GET("/:id", projectHandler.GetProject)
		privateProjectRoutes.GET("", projectHandler.GetAllProjects)
		privateProjectRoutes.POST("/comments", projectHandler.CreateComment)
		privateProjectRoutes.DELETE("/comments/:id", projectHandler.DeleteComment)
		privateProjectRoutes.POST("/reviews", projectHandler.CreateReview)
	}
}

func RegisterAdminProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	adminProjectRoutes := r.Group("/api/admin/projects")
	{
		adminProjectRoutes.DELETE("/:id", projectHandler.DeleteProject)
		adminProjectRoutes.PATCH("/comments/:id/moderate", projectHandler.ModerateComment)
	}
}

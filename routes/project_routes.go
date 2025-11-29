package routes

import (
	"github.com/aruncs31s/esdcprojectmodule/handler"
	"github.com/gin-gonic/gin"
)

func RegisterPublicProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	publicProjectRoutes := r.Group("/api/public/projects")
	{
		publicProjectRoutes.GET("", projectHandler.GetAllProjects)
		publicProjectRoutes.GET("/trending", projectHandler.GetTrendingProjects)
		publicProjectRoutes.GET("/analytics/platform", projectHandler.GetPlatformAnalytics)
		publicProjectRoutes.GET("/templates", projectHandler.GetTemplates)
		publicProjectRoutes.GET("/:id", projectHandler.GetProject)
		publicProjectRoutes.GET("/:id/similar", projectHandler.GetSimilarProjects)
	}
}

func RegisterPrivateProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	privateProjectRoutes := r.Group("/api/projects")
	{
		// Project management
		privateProjectRoutes.POST("", projectHandler.CreateProject)
		privateProjectRoutes.POST("/:id/toggle-like", projectHandler.ToggleLikeProject)
		privateProjectRoutes.GET("/:id", projectHandler.GetProject)
		privateProjectRoutes.GET("", projectHandler.GetAllProjects)

		// Trending & Recommendations
		privateProjectRoutes.GET("/recommendations", projectHandler.GetRecommendedProjects)
		privateProjectRoutes.GET("/:id/analytics", projectHandler.GetProjectAnalytics)

		// Templates
		privateProjectRoutes.POST("/templates", projectHandler.CreateTemplate)
		privateProjectRoutes.GET("/templates/user", projectHandler.GetUserTemplates)
		privateProjectRoutes.GET("/templates/:id", projectHandler.GetTemplate)
		privateProjectRoutes.DELETE("/templates/:id", projectHandler.DeleteTemplate)
		privateProjectRoutes.POST("/from-template", projectHandler.CreateProjectFromTemplate)

		// Notifications
		privateProjectRoutes.GET("/notifications", projectHandler.GetNotifications)
		privateProjectRoutes.POST("/notifications/:id/read", projectHandler.MarkNotificationAsRead)
		privateProjectRoutes.POST("/notifications/read-all", projectHandler.MarkAllNotificationsAsRead)
		privateProjectRoutes.DELETE("/notifications/:id", projectHandler.DeleteNotification)

		// Export
		privateProjectRoutes.GET("/:id/export", projectHandler.ExportProject)
		privateProjectRoutes.GET("/portfolio/export", projectHandler.ExportPortfolio)
	}
}

func RegisterAdminProjectRoutes(r *gin.Engine, projectHandler handler.ProjectHandler) {
	// Admin routes for future use
	adminProjectRoutes := r.Group("/api/admin/projects")
	{
		_ = adminProjectRoutes // Placeholder for future admin functionality
	}
}

package service

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
)

type ProjectService interface {
	CreateProject(user string, project dto.ProjectCreation) (*commonModules.Project, error)
	GetProject(id uint, user string) (*dto.ProjectResponse, error)
	GetUserProjects(limit, offset int, username string) ([]*dto.ProjectResponse, error)
	ToggleLikeProject(username string, projectID uint) (bool, error)
	FilterProjects(filter dto.ProjectFilter) (*dto.ProjectListResponse, error)
	UpdateProject(projectID uint, username string, updates dto.ProjectUpdate) error
	DeleteProject(projectID uint, username string, isAdmin bool) error
	GetProjectStats(projectID uint) (*dto.ProjectStats, error)
	CreateComment(username string, comment dto.CommentCreate) (*dto.CommentResponse, error)
	GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error)
	CreateReview(username string, review dto.ReviewCreate) (*dto.ReviewResponse, error)
	GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error)
	DeleteComment(commentID uint, username string, isAdmin bool) error
	ModerateComment(commentID uint, status string) error
	// Trending & Recommendations
	GetTrendingProjects(limit, offset int, days int) ([]dto.TrendingProject, error)
	GetRecommendedProjects(username string, limit, offset int) ([]dto.RecommendedProject, error)
	GetSimilarProjects(projectID uint, limit, offset int) ([]dto.SimilarProject, error)
	// Analytics
	GetProjectAnalytics(projectID uint, username string, days int) (*dto.ProjectAnalytics, error)
	GetPlatformAnalytics(days int) (*dto.PlatformAnalytics, error)
	// Templates
	CreateTemplate(username string, template dto.TemplateCreate) (*dto.ProjectTemplate, error)
	GetTemplates(limit, offset int, isPublic bool) ([]dto.TemplateListItem, error)
	GetUserTemplates(username string, limit, offset int) ([]dto.TemplateListItem, error)
	GetTemplate(templateID uint) (*dto.ProjectTemplate, error)
	DeleteTemplate(templateID uint, username string) error
	CreateProjectFromTemplate(username string, request dto.CreateFromTemplate) (*dto.ProjectResponse, error)
	// Notifications
	GetNotifications(username string, limit, offset int) ([]dto.NotificationResponse, error)
	MarkNotificationAsRead(notificationID uint, username string) error
	MarkAllAsRead(username string) error
	DeleteNotification(notificationID uint, username string) error
	CreateNotification(userID uint, notification dto.NotificationCreate) error
}

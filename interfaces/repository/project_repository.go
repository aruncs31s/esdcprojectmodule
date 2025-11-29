package repository

import (
	model "github.com/aruncs31s/esdcmodels"
)

type ProjectRepository interface {
	ProjectRepositoryReader
	ProjectRepositoryWriter
	ProjectRepositoryMixed
}

// CreateProject is used to create a new project with the provided details.
//
// Currently Used In create Project Modal // Remove this comment later
// Params:
//  - user: string - The username of the creator.
// - project: dto.ProjectCreation - The details of the project to be created.
// - Returns:
// - *commonModules.Project - The created project.
// - error - An error if the creation fails.

type ProjectRepositoryReader interface {
	GetPublicProjects(limit, offset int) ([]model.Project, error)
	GetUserProjects(userID uint, limit, offset int) ([]model.Project, error)
	GetEssentialInfo(limit, offset int) (*[]model.Project, error)
	GetByID(id uint) (model.Project, error)
	GetProjectsCount() (int, error)
	IsLiked(userID uint, projectID uint) (bool, error)
	FilterProjects(category, status, visibility, search string, technologies []string, limit, offset int) ([]model.Project, int, error)
	GetProjectStats(projectID uint) (*model.ProjectStats, error)
	GetComments(projectID uint, limit, offset int) ([]model.Comment, error)
	GetReviews(projectID uint, limit, offset int) ([]model.Review, error)
	// Trending & Recommendations
	GetTrendingProjects(limit, offset int, days int) ([]model.Project, error)
	GetRecommendedProjects(userID uint, limit, offset int) ([]model.Project, error)
	GetSimilarProjects(projectID uint, limit, offset int) ([]model.Project, error)
	// Analytics
	GetProjectAnalytics(projectID uint, days int) (*model.ProjectAnalytics, error)
	GetPlatformAnalytics(days int) (*model.PlatformAnalytics, error)
	GetTrendingTechnologies(limit int) ([]model.TrendingTech, error)
	GetTrendingTags(limit int) ([]model.TrendingTag, error)
	// Templates
	GetTemplateByID(templateID uint) (*model.ProjectTemplate, error)
	GetUserTemplates(userID uint, limit, offset int) ([]model.ProjectTemplate, error)
	GetPublicTemplates(limit, offset int) ([]model.ProjectTemplate, error)
	// Notifications
	GetUserNotifications(userID uint, limit, offset int) ([]model.Notification, error)
	GetUnreadNotificationCount(userID uint) (int, error)
}
type ProjectRepositoryMixed interface {
	// FindOrCreateTag finds a tag by name or creates it if it doesn't exist.
	//
	// Params:
	//   - name: string - The name of the tag to find or create.
	//
	// Returns:
	//   - *commonModules.Tag: A pointer to the Tag object.
	//   - error: An error object if any error occurs during the database operation.
	FindOrCreateTag(name string) (*model.Tag, error)
	// FindOrCreateTechnology finds a technology by name or creates it if it doesn't exist.
	//
	// Params:
	//   - name: string - The name of the technology to find or create.
	//
	// Returns:
	//   - *commonModels.Technologies: A pointer to the Technologies object.
	//   - error: An error object if any error occurs during the database operation.
	FindOrCreateTechnology(name string) (*model.Technologies, error)
}
type ProjectRepositoryWriter interface {
	Create(project *model.Project) error
	Update(projectID uint, updates map[string]interface{}) error
	Delete(projectID uint, soft bool) error
	LikeProject(userID uint, projectID uint) error
	UnlikeProject(userID uint, projectID uint) error
	IncrementViewCount(projectID uint) error
	CreateComment(comment *model.Comment) error
	CreateReview(review *model.Review) error
	DeleteComment(commentID uint) error
	UpdateCommentStatus(commentID uint, status string) error
	// Templates
	CreateTemplate(template *model.ProjectTemplate) error
	DeleteTemplate(templateID uint) error
	IncrementTemplateUsage(templateID uint) error
	// Notifications
	CreateNotification(notification *model.Notification) error
	MarkNotificationAsRead(notificationID uint) error
	MarkAllNotificationsAsRead(userID uint) error
	DeleteNotification(notificationID uint) error
}

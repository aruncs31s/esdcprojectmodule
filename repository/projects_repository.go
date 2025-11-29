package repository

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/repository"

	"gorm.io/gorm"
)

type projectRepositoryMixed struct {
	db *gorm.DB
}

type projectRepositoryReader struct {
	db *gorm.DB
}
type projectRepositoryWriter struct {
	db *gorm.DB
}
type projectRepository struct {
	reader repository.ProjectRepositoryReader
	writer repository.ProjectRepositoryWriter
	mixed  repository.ProjectRepositoryMixed
}

func newProjectRepositoryReader(db *gorm.DB) repository.ProjectRepositoryReader {
	return &projectRepositoryReader{db: db}
}
func newProjectRepositoryWriter(db *gorm.DB) repository.ProjectRepositoryWriter {
	return &projectRepositoryWriter{db: db}
}
func newProjectRepositoryMixed(db *gorm.DB) repository.ProjectRepositoryMixed {
	return &projectRepositoryMixed{db: db}
}

func NewProjectRepository(db *gorm.DB) repository.ProjectRepository {
	return &projectRepository{
		reader: newProjectRepositoryReader(db),
		writer: newProjectRepositoryWriter(db),
		mixed:  newProjectRepositoryMixed(db),
	}
}

func (r *projectRepository) GetPublicProjects(limit, offset int) ([]commonModules.Project, error) {
	return r.reader.GetPublicProjects(limit, offset)
}

func (r *projectRepository) GetUserProjects(userID uint, limit, offset int) ([]commonModules.Project, error) {
	return r.reader.GetUserProjects(userID, limit, offset)
}

func (r *projectRepository) GetByID(id uint) (commonModules.Project, error) {
	return r.reader.GetByID(id)
}

func (r *projectRepository) GetEssentialInfo(limit, offset int) (*[]commonModules.Project, error) {
	return r.reader.GetEssentialInfo(limit, offset)
}

func (r *projectRepository) GetProjectsCount() (int, error) {
	return r.reader.GetProjectsCount()
}

func (r *projectRepository) IsLiked(userID uint, projectID uint) (bool, error) {
	return r.reader.IsLiked(userID, projectID)
}

func (r *projectRepository) Create(project *commonModules.Project) error {
	return r.writer.Create(project)
}

func (r *projectRepository) LikeProject(userID uint, projectID uint) error {
	return r.writer.LikeProject(userID, projectID)
}

func (r *projectRepository) UnlikeProject(userID uint, projectID uint) error {
	return r.writer.UnlikeProject(userID, projectID)
}

func (r *projectRepository) FindOrCreateTag(name string) (*commonModules.Tag, error) {
	return r.mixed.FindOrCreateTag(name)
}

func (r *projectRepository) FindOrCreateTechnology(name string) (*commonModules.Technologies, error) {
	return r.mixed.FindOrCreateTechnology(name)
}

func (r *projectRepository) FilterProjects(category, status, visibility, search string, technologies []string, limit, offset int) ([]commonModules.Project, int, error) {
	return r.reader.FilterProjects(category, status, visibility, search, technologies, limit, offset)
}

func (r *projectRepository) GetProjectStats(projectID uint) (*dto.ProjectStats, error) {
	return r.reader.GetProjectStats(projectID)
}

func (r *projectRepository) GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error) {
	return r.reader.GetComments(projectID, limit, offset)
}

func (r *projectRepository) GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error) {
	return r.reader.GetReviews(projectID, limit, offset)
}

func (r *projectRepository) Update(projectID uint, updates map[string]interface{}) error {
	return r.writer.Update(projectID, updates)
}

func (r *projectRepository) Delete(projectID uint, soft bool) error {
	return r.writer.Delete(projectID, soft)
}

func (r *projectRepository) IncrementViewCount(projectID uint) error {
	return r.writer.IncrementViewCount(projectID)
}

func (r *projectRepository) CreateComment(projectID, userID uint, content string) (*dto.CommentResponse, error) {
	return r.writer.CreateComment(projectID, userID, content)
}

func (r *projectRepository) CreateReview(projectID, userID uint, rating float64, comment string) (*dto.ReviewResponse, error) {
	return r.writer.CreateReview(projectID, userID, rating, comment)
}

func (r *projectRepository) DeleteComment(commentID uint) error {
	return r.writer.DeleteComment(commentID)
}

func (r *projectRepository) UpdateCommentStatus(commentID uint, status string) error {
	return r.writer.UpdateCommentStatus(commentID, status)
}

func (r *projectRepositoryReader) GetPublicProjects(limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("visibility = ?", 0).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepositoryReader) GetUserProjects(userID uint, limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("created_by = ? OR visibility = 0", userID).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepositoryReader) GetByID(id uint) (commonModules.Project, error) {
	var project commonModules.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		First(&project, id).Error; err != nil {
		return commonModules.Project{}, err
	}
	if project.IsPrivate() {
		return commonModules.Project{}, gorm.ErrRecordNotFound
	}
	return project, nil
}

func (r *projectRepositoryReader) GetProjectsCount() (int, error) {
	var count int64
	result := r.db.Model(&commonModules.Project{}).Count(&count)
	return int(count), result.Error
}

func (r *projectRepositoryReader) IsLiked(userID uint, projectID uint) (bool, error) {
	var count int64
	err := r.db.Table("project_likes").Where("user_id = ? AND project_id = ?", userID, projectID).Count(&count).Error
	return count > 0, err
}

func (r *projectRepositoryWriter) Create(project *commonModules.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}
	return nil
}

func (r *projectRepositoryWriter) LikeProject(userID uint, projectID uint) error {
	// Add to association
	var user commonModules.User
	var project commonModules.Project
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	if err := r.db.First(&project, projectID).Error; err != nil {
		return err
	}
	if err := r.db.Model(&user).Association("LikedProjects").Append(&project); err != nil {
		return err
	}
	// Update likes count
	return r.db.Model(&commonModules.Project{}).Where("id = ?", projectID).Update("likes", gorm.Expr("likes + ?", 1)).Error
}

func (r *projectRepositoryWriter) UnlikeProject(userID uint, projectID uint) error {
	// Remove from association
	var user commonModules.User
	var project commonModules.Project
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}
	if err := r.db.First(&project, projectID).Error; err != nil {
		return err
	}
	if err := r.db.Model(&user).Association("LikedProjects").Delete(&project); err != nil {
		return err
	}
	// Update likes count
	return r.db.Model(&commonModules.Project{}).Where("id = ?", projectID).Update("likes", gorm.Expr("likes - ?", 1)).Error
}

func (r *projectRepositoryMixed) FindOrCreateTag(name string) (*commonModules.Tag, error) {
	var tag commonModules.Tag
	if err := r.db.Where("name = ?", name).FirstOrCreate(&tag, commonModules.Tag{Name: name}).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *projectRepositoryMixed) FindOrCreateTechnology(name string) (*commonModules.Technologies, error) {
	var tech commonModules.Technologies
	if err := r.db.Where("name = ?", name).FirstOrCreate(&tech, commonModules.Technologies{Name: name}).Error; err != nil {
		return nil, err
	}
	return &tech, nil
}

func (r *projectRepositoryReader) FilterProjects(category, status, visibility, search string, technologies []string, limit, offset int) ([]commonModules.Project, int, error) {
	query := r.db.Model(&commonModules.Project{}).Preload("Contributors").Preload("Creator").Preload("Tags").Preload("Technologies")

	if category != "" {
		query = query.Where("category = ?", category)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if visibility != "" {
		query = query.Where("visibility = ?", visibility)
	}
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if len(technologies) > 0 {
		query = query.Joins("JOIN project_technologies ON projects.id = project_technologies.project_id").
			Joins("JOIN technologies ON project_technologies.technology_id = technologies.id").
			Where("technologies.name IN ?", technologies)
	}

	var total int64
	query.Count(&total)

	var projects []commonModules.Project
	err := query.Limit(limit).Offset(offset).Find(&projects).Error
	return projects, int(total), err
}

func (r *projectRepositoryReader) GetProjectStats(projectID uint) (*dto.ProjectStats, error) {
	var stats dto.ProjectStats
	var project commonModules.Project

	if err := r.db.First(&project, projectID).Error; err != nil {
		return nil, err
	}

	var commentCount, reviewCount int64
	var avgRating float64

	// Count comments from project_comments table
	r.db.Table("project_comments").Where("project_id = ?", projectID).Count(&commentCount)
	// Count reviews from project_reviews table
	r.db.Table("project_reviews").Where("project_id = ?", projectID).Count(&reviewCount)
	// Get average rating
	r.db.Table("project_reviews").Where("project_id = ?", projectID).Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	stats.ViewCount = project.Views
	stats.LikeCount = project.Likes
	stats.CommentCount = int(commentCount)
	stats.ReviewCount = int(reviewCount)
	stats.AverageRating = avgRating

	return &stats, nil
}

func (r *projectRepositoryReader) GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error) {
	var comments []dto.CommentResponse
	// Query directly from project_comments table and join with users
	err := r.db.Table("project_comments").
		Select("project_comments.id, project_comments.project_id, project_comments.content, project_comments.created_at, users.id as user_id, users.name as user_name, users.email as user_email, users.image as user_image").
		Joins("LEFT JOIN users ON project_comments.user_id = users.id").
		Where("project_comments.project_id = ?", projectID).
		Limit(limit).Offset(offset).
		Order("project_comments.created_at DESC").
		Scan(&comments).Error
	return comments, err
}

func (r *projectRepositoryReader) GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error) {
	var reviews []dto.ReviewResponse
	// Query directly from project_reviews table and join with users
	err := r.db.Table("project_reviews").
		Select("project_reviews.id, project_reviews.project_id, project_reviews.rating, project_reviews.comment, project_reviews.created_at, users.id as user_id, users.name as user_name, users.email as user_email, users.image as user_image").
		Joins("LEFT JOIN users ON project_reviews.user_id = users.id").
		Where("project_reviews.project_id = ?", projectID).
		Limit(limit).Offset(offset).
		Order("project_reviews.created_at DESC").
		Scan(&reviews).Error
	return reviews, err
}

func (r *projectRepositoryWriter) Update(projectID uint, updates map[string]interface{}) error {
	return r.db.Model(&commonModules.Project{}).Where("id = ?", projectID).Updates(updates).Error
}

func (r *projectRepositoryWriter) Delete(projectID uint, soft bool) error {
	if soft {
		return r.db.Delete(&commonModules.Project{}, projectID).Error
	}
	return r.db.Unscoped().Delete(&commonModules.Project{}, projectID).Error
}

func (r *projectRepositoryWriter) IncrementViewCount(projectID uint) error {
	return r.db.Model(&commonModules.Project{}).Where("id = ?", projectID).Update("views", gorm.Expr("views + ?", 1)).Error
}

func (r *projectRepositoryWriter) CreateComment(projectID, userID uint, content string) (*dto.CommentResponse, error) {
	// Insert into project_comments table
	result := r.db.Exec("INSERT INTO project_comments (project_id, user_id, content, status, created_at, updated_at) VALUES (?, ?, ?, 'approved', NOW(), NOW())", projectID, userID, content)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return the created comment response
	var response dto.CommentResponse
	r.db.Table("project_comments").
		Select("project_comments.id, project_comments.project_id, project_comments.content, project_comments.created_at").
		Where("project_comments.project_id = ? AND project_comments.user_id = ?", projectID, userID).
		Order("project_comments.id DESC").
		Limit(1).
		Scan(&response)
	return &response, nil
}

func (r *projectRepositoryWriter) CreateReview(projectID, userID uint, rating float64, comment string) (*dto.ReviewResponse, error) {
	// Insert into project_reviews table
	result := r.db.Exec("INSERT INTO project_reviews (project_id, user_id, rating, comment, created_at) VALUES (?, ?, ?, ?, NOW())", projectID, userID, rating, comment)
	if result.Error != nil {
		return nil, result.Error
	}
	// Return the created review response
	var response dto.ReviewResponse
	r.db.Table("project_reviews").
		Select("project_reviews.id, project_reviews.project_id, project_reviews.rating, project_reviews.comment, project_reviews.created_at").
		Where("project_reviews.project_id = ? AND project_reviews.user_id = ?", projectID, userID).
		Order("project_reviews.id DESC").
		Limit(1).
		Scan(&response)
	return &response, nil
}

func (r *projectRepositoryWriter) DeleteComment(commentID uint) error {
	return r.db.Exec("DELETE FROM project_comments WHERE id = ?", commentID).Error
}

func (r *projectRepositoryWriter) UpdateCommentStatus(commentID uint, status string) error {
	return r.db.Exec("UPDATE project_comments SET status = ? WHERE id = ?", status, commentID).Error
}

func (r *projectRepositoryReader) GetEssentialInfo(limit, offset int) (*[]commonModules.Project, error) {
	var projects []commonModules.Project
	err := r.db.Limit(limit).Offset(offset).Find(&projects).Error
	return &projects, err
}

// ========== DELEGATION METHODS FOR TRENDING/RECOMMENDATIONS/ANALYTICS ==========

func (r *projectRepository) GetTrendingProjects(limit, offset int, days int) ([]commonModules.Project, error) {
	return r.reader.GetTrendingProjects(limit, offset, days)
}

func (r *projectRepository) GetRecommendedProjects(userID uint, limit, offset int) ([]commonModules.Project, error) {
	return r.reader.GetRecommendedProjects(userID, limit, offset)
}

func (r *projectRepository) GetSimilarProjects(projectID uint, limit, offset int) ([]commonModules.Project, error) {
	return r.reader.GetSimilarProjects(projectID, limit, offset)
}

func (r *projectRepository) GetProjectAnalytics(projectID uint, days int) (*dto.ProjectStats, error) {
	return r.reader.GetProjectAnalytics(projectID, days)
}

func (r *projectRepository) GetPlatformAnalytics(days int) (*dto.PlatformAnalytics, error) {
	return r.reader.GetPlatformAnalytics(days)
}

func (r *projectRepository) GetTrendingTechnologies(limit int) ([]commonModules.TrendingTech, error) {
	return r.reader.GetTrendingTechnologies(limit)
}

func (r *projectRepository) GetTrendingTags(limit int) ([]commonModules.TrendingTag, error) {
	return r.reader.GetTrendingTags(limit)
}

// ========== DELEGATION METHODS FOR TEMPLATES ==========

func (r *projectRepository) GetTemplateByID(templateID uint) (*commonModules.ProjectTemplate, error) {
	return r.reader.GetTemplateByID(templateID)
}

func (r *projectRepository) GetUserTemplates(userID uint, limit, offset int) ([]commonModules.ProjectTemplate, error) {
	return r.reader.GetUserTemplates(userID, limit, offset)
}

func (r *projectRepository) GetPublicTemplates(limit, offset int) ([]commonModules.ProjectTemplate, error) {
	return r.reader.GetPublicTemplates(limit, offset)
}

func (r *projectRepository) CreateTemplate(template *commonModules.ProjectTemplate) error {
	return r.writer.CreateTemplate(template)
}

func (r *projectRepository) DeleteTemplate(templateID uint) error {
	return r.writer.DeleteTemplate(templateID)
}

func (r *projectRepository) IncrementTemplateUsage(templateID uint) error {
	return r.writer.IncrementTemplateUsage(templateID)
}

// ========== DELEGATION METHODS FOR NOTIFICATIONS ==========

func (r *projectRepository) GetUserNotifications(userID uint, limit, offset int) ([]commonModules.Notification, error) {
	return r.reader.GetUserNotifications(userID, limit, offset)
}

func (r *projectRepository) GetUnreadNotificationCount(userID uint) (int, error) {
	return r.reader.GetUnreadNotificationCount(userID)
}

func (r *projectRepository) CreateNotification(notification *commonModules.Notification) error {
	return r.writer.CreateNotification(notification)
}

func (r *projectRepository) MarkNotificationAsRead(notificationID uint) error {
	return r.writer.MarkNotificationAsRead(notificationID)
}

func (r *projectRepository) MarkAllNotificationsAsRead(userID uint) error {
	return r.writer.MarkAllNotificationsAsRead(userID)
}

func (r *projectRepository) DeleteNotification(notificationID uint) error {
	return r.writer.DeleteNotification(notificationID)
}

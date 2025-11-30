package repository

import (
	model "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
	repository "github.com/aruncs31s/esdcprojectmodule/interfaces/repository"
	"gorm.io/gorm"
)

type publicProjectRepository struct {
	db *gorm.DB
}

func NewPublicProjectRepository(db *gorm.DB) repository.PublicProjectRepository {
	return &publicProjectRepository{
		db: db,
	}
}
func (r *publicProjectRepository) GetAllProjects(limit, offset int) (*[]model.Project, error) {
	var projects []model.Project
	// TODO: Move the visibility to model.
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Preload("LikedBy").
		Preload("ViewedBy").
		Preload("Comments").
		Preload("Reviews").
		Where("visibility = ?", 0).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return &projects, nil
}
func (r *publicProjectRepository) GetUserProjects(userID uint, limit, offset int) (*[]model.Project, error) {
	var projects []model.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Tags").
		Preload("Technologies").
		Where("created_by = ? OR visibility = 0", userID).
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return &projects, nil
}
func (r *publicProjectRepository) GetProject(id uint) (*model.Project, error) {
	var project model.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Preload("LikedBy").
		Preload("ViewedBy").
		Preload("Comments").
		Preload("Reviews").
		First(&project, id).Error; err != nil {
		return nil, err
	}
	if project.IsPrivate() {
		return nil, gorm.ErrRecordNotFound
	}
	return &project, nil
}

func (r *publicProjectRepository) GetProjectsCount() (int, error) {
	var count int64
	if err := r.db.Where("visibility = ?", 0).Model(&model.Project{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *publicProjectRepository) GetProjectStats(projectID uint) (*dto.ProjectStats, error) {
	var stats dto.ProjectStats
	var project model.Project

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

func (r *publicProjectRepository) GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error) {
	var comments []dto.CommentResponse
	err := r.db.Table("project_comments").
		Select("project_comments.id, project_comments.project_id, project_comments.content, project_comments.created_at, users.id as user_id, users.name as user_name, users.email as user_email, users.image as user_image").
		Joins("LEFT JOIN users ON project_comments.user_id = users.id").
		Where("project_comments.project_id = ?", projectID).
		Limit(limit).Offset(offset).
		Order("project_comments.created_at DESC").
		Scan(&comments).Error
	return comments, err
}

func (r *publicProjectRepository) GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error) {
	var reviews []dto.ReviewResponse
	err := r.db.Table("project_reviews").
		Select("project_reviews.id, project_reviews.project_id, project_reviews.rating, project_reviews.comment, project_reviews.created_at, users.id as user_id, users.name as user_name, users.email as user_email, users.image as user_image").
		Joins("LEFT JOIN users ON project_reviews.user_id = users.id").
		Where("project_reviews.project_id = ?", projectID).
		Limit(limit).Offset(offset).
		Order("project_reviews.created_at DESC").
		Scan(&reviews).Error
	return reviews, err
}

func (r *publicProjectRepository) GetTrendingProjects(limit, offset, days int) ([]model.Project, error) {
	var projects []model.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("visibility = ? AND created_at >= DATE_SUB(NOW(), INTERVAL ? DAY)", 0, days).
		Order("likes DESC, views DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *publicProjectRepository) GetSimilarProjects(projectID uint, limit, offset int) ([]model.Project, error) {
	var referenceProject model.Project
	if err := r.db.First(&referenceProject, projectID).Error; err != nil {
		return nil, err
	}

	var projects []model.Project
	if err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("visibility = ? AND id != ? AND category = ?", 0, projectID, referenceProject.Category).
		Order("likes DESC, views DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *publicProjectRepository) GetPublicTemplates(limit, offset int) ([]model.ProjectTemplate, error) {
	var templates []model.ProjectTemplate
	if err := r.db.
		Preload("Creator").
		Preload("Technologies").
		Preload("Tags").
		Where("is_public = ?", true).
		Limit(limit).
		Offset(offset).
		Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *publicProjectRepository) GetTemplateByID(templateID uint) (*model.ProjectTemplate, error) {
	var template model.ProjectTemplate
	if err := r.db.
		Preload("Creator").
		Preload("Technologies").
		Preload("Tags").
		First(&template, templateID).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

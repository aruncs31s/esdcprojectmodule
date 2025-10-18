package repository

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	model "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetAll(limit, offset int) ([]commonModules.Project, error)
	GetPublicProjects(limit, offset int) ([]commonModules.Project, error)
	GetUserProjects(userID uint, limit, offset int) ([]commonModules.Project, error)
	GetEssentialInfo(limit, offset int) ([]dto.ProjectsEssentialInfo, error)
	Create(project *commonModules.Project) error
	GetByID(id int) (commonModules.Project, error)
	// Update(id int, project model.Project) (model.Project, error)
	// Delete(id int) error
	GetProjectsCount() (int, error)
	FindOrCreateTag(name string) (*commonModules.Tag, error)
	FindOrCreateTechnology(name string) (*commonModules.Technologies, error)
	IsLiked(userID uint, projectID int) (bool, error)
	LikeProject(userID uint, projectID int) error
	UnlikeProject(userID uint, projectID int) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// For Admin Pannel
func (r *projectRepository) GetEssentialInfo(limit, offset int) ([]dto.ProjectsEssentialInfo, error) {
	var projects []dto.ProjectsEssentialInfo
	query := r.db.
		Model(&commonModules.Project{}).
		Select("projects.id, projects.title, users.name as created_by, projects.status, projects.created_at, projects.updated_at").
		Joins("LEFT JOIN users ON projects.created_by = users.id").
		Limit(limit).Offset(offset)
	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetAll(limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project
	if err := r.db.Preload("Contributors").Preload("Creator").Preload("Tags").Preload("Technologies").Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetPublicProjects(limit, offset int) ([]model.Project, error) {
	var projects []commonModules.Project
	if err := r.db.Preload("Contributors").Preload("Creator").Preload("Tags").Preload("Technologies").Where("visibility = ?", 0).Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *projectRepository) GetUserProjects(userID uint, limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project
	if err := r.db.Preload("Contributors").Preload("Creator").Preload("Tags").Preload("Technologies").Where("created_by = ? OR visibility = 0", userID).Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}
func (r *projectRepository) GetByID(id int) (commonModules.Project, error) {
	var project commonModules.Project
	if err := r.db.Preload("Contributors").Preload("Creator").Preload("Tags").Preload("Technologies").First(&project, id).Error; err != nil {
		return commonModules.Project{}, err
	}
	return project, nil
}
func (r *projectRepository) Create(project *commonModules.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}
	return nil
}
func (r *projectRepository) GetProjectsCount() (int, error) {
	var count int64
	result := r.db.Model(&commonModules.Project{}).Count(&count)
	return int(count), result.Error
}

func (r *projectRepository) FindOrCreateTag(name string) (*commonModules.Tag, error) {
	var tag commonModules.Tag
	if err := r.db.Where("name = ?", name).FirstOrCreate(&tag, commonModules.Tag{Name: name}).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *projectRepository) FindOrCreateTechnology(name string) (*commonModules.Technologies, error) {
	var tech commonModules.Technologies
	if err := r.db.Where("name = ?", name).FirstOrCreate(&tech, commonModules.Technologies{Name: name}).Error; err != nil {
		return nil, err
	}
	return &tech, nil
}

func (r *projectRepository) IsLiked(userID uint, projectID int) (bool, error) {
	var count int64
	err := r.db.Table("project_likes").Where("user_id = ? AND project_id = ?", userID, projectID).Count(&count).Error
	return count > 0, err
}

func (r *projectRepository) LikeProject(userID uint, projectID int) error {
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

func (r *projectRepository) UnlikeProject(userID uint, projectID int) error {
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

package repository

import (
	model "github.com/aruncs31s/esdcmodels"
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

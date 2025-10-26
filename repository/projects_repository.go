package repository

import (
	commonModules "github.com/aruncs31s/esdcmodels"
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

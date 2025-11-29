package repository

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"gorm.io/gorm"
)

// GetTemplateByID retrieves a template by ID with all relationships
func (r *projectRepositoryReader) GetTemplateByID(templateID uint) (*commonModules.ProjectTemplate, error) {
	var template commonModules.ProjectTemplate
	err := r.db.
		Preload("Creator").
		Preload("Technologies").
		Preload("Tags").
		First(&template, templateID).Error
	return &template, err
}

// GetUserTemplates retrieves all templates created by a specific user
func (r *projectRepositoryReader) GetUserTemplates(userID uint, limit, offset int) ([]commonModules.ProjectTemplate, error) {
	var templates []commonModules.ProjectTemplate
	err := r.db.
		Preload("Creator").
		Preload("Technologies").
		Preload("Tags").
		Where("created_by = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&templates).Error
	return templates, err
}

// GetPublicTemplates retrieves all public templates
func (r *projectRepositoryReader) GetPublicTemplates(limit, offset int) ([]commonModules.ProjectTemplate, error) {
	var templates []commonModules.ProjectTemplate
	err := r.db.
		Preload("Creator").
		Preload("Technologies").
		Preload("Tags").
		Where("is_public = ?", true).
		Order("usage_count DESC, created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&templates).Error
	return templates, err
}

// CreateTemplate saves a new project template
func (w *projectRepositoryWriter) CreateTemplate(template *commonModules.ProjectTemplate) error {
	return w.db.Create(template).Error
}

// DeleteTemplate removes a template
func (w *projectRepositoryWriter) DeleteTemplate(templateID uint) error {
	return w.db.Delete(&commonModules.ProjectTemplate{}, templateID).Error
}

// IncrementTemplateUsage increases the usage count of a template
func (w *projectRepositoryWriter) IncrementTemplateUsage(templateID uint) error {
	return w.db.Model(&commonModules.ProjectTemplate{}).
		Where("id = ?", templateID).
		Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error
}

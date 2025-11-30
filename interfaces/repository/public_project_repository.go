package repository

import (
	model "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
)

type PublicProjectRepository interface {
	// GetAllProjects retrieves all public projects.
	//
	// Params:
	//  - limit: int - The maximum number of projects to retrieve.
	//  - offset: int - The number of projects to skip before starting to collect the result set.
	// Returns:
	//  - []model.Project - A slice of Project models representing the public projects.
	//  - error - An error if the retrieval fails.
	GetAllProjects(limit, offset int) (*[]model.Project, error)

	// GetUserProjects retrieves all public projects for a specific user.
	//
	// Params:
	//  - user: uint - The user ID of the project creator.
	//  - limit: int - The maximum number of projects to retrieve.
	//  - offset: int - The number of projects to skip before starting to collect the result set.
	// Returns:
	//  - []model.Project - A slice of Project models representing the user's public projects.
	//  - error - An error if the retrieval fails.
	GetUserProjects(user uint, limit, offset int) (*[]model.Project, error)

	// GetProject retrieves a single project by ID.
	//
	// Params:
	//  - id: uint - The ID of the project to retrieve.
	// Returns:
	//  - *model.Project - The project model.
	//  - error - An error if the retrieval fails.
	GetProject(id uint) (*model.Project, error)

	// GetProjectsCount retrieves the total count of public projects.
	//
	// Returns:
	//  - int - The total count of public projects.
	//  - error - An error if the retrieval fails.
	GetProjectsCount() (int, error)

	// GetProjectStats retrieves statistics for a specific project.
	//
	// Params:
	//  - projectID: uint - The ID of the project.
	// Returns:
	//  - *model.ProjectStats - The project statistics.
	//  - error - An error if the retrieval fails.
	GetProjectStats(projectID uint) (*dto.ProjectStats, error)

	// GetComments retrieves comments for a specific project.
	//
	// Params:
	//  - projectID: uint - The ID of the project.
	//  - limit: int - The maximum number of comments to retrieve.
	//  - offset: int - The number of comments to skip before starting to collect the result set.
	// Returns:
	//  - []dto.CommentResponse - A slice of comment responses for the project.
	//  - error - An error if the retrieval fails.
	GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error)

	// GetReviews retrieves reviews for a specific project.
	//
	// Params:
	//  - projectID: uint - The ID of the project.
	//  - limit: int - The maximum number of reviews to retrieve.
	//  - offset: int - The number of reviews to skip before starting to collect the result set.
	// Returns:
	//  - []dto.ReviewResponse - A slice of review responses for the project.
	//  - error - An error if the retrieval fails.
	GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error)

	// GetTrendingProjects retrieves trending projects.
	//
	// Params:
	//  - limit: int - The maximum number of projects to retrieve.
	//  - offset: int - The number of projects to skip before starting to collect the result set.
	//  - days: int - The number of days to consider for trending calculation.
	// Returns:
	//  - []model.Project - A slice of trending projects.
	//  - error - An error if the retrieval fails.
	GetTrendingProjects(limit, offset int, days int) ([]model.Project, error)

	// GetSimilarProjects retrieves projects similar to a specific project.
	//
	// Params:
	//  - projectID: uint - The ID of the project to find similar projects for.
	//  - limit: int - The maximum number of projects to retrieve.
	//  - offset: int - The number of projects to skip before starting to collect the result set.
	// Returns:
	//  - []model.Project - A slice of similar projects.
	//  - error - An error if the retrieval fails.
	GetSimilarProjects(projectID uint, limit, offset int) ([]model.Project, error)

	// GetPublicTemplates retrieves public project templates.
	//
	// Params:
	//  - limit: int - The maximum number of templates to retrieve.
	//  - offset: int - The number of templates to skip before starting to collect the result set.
	// Returns:
	//  - []model.ProjectTemplate - A slice of public project templates.
	//  - error - An error if the retrieval fails.
	GetPublicTemplates(limit, offset int) ([]model.ProjectTemplate, error)

	// GetTemplateByID retrieves a specific project template by ID.
	//
	// Params:
	//  - templateID: uint - The ID of the template to retrieve.
	// Returns:
	//  - *model.ProjectTemplate - The project template model.
	//  - error - An error if the retrieval fails.
	GetTemplateByID(templateID uint) (*model.ProjectTemplate, error)
}

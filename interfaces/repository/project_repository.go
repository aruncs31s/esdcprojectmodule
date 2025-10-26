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
	// Used By Public Routes and Homepage
	// GetPublicProjects retrieves public projects with pagination.
	//
	// It fetches projects that are marked as public (visibility = 0).
	//
	// Params:
	//   - limit: int - The maximum number of projects to retrieve.
	//   - offset: int - The number of projects to skip before starting to collect the result set.
	//
	// Returns:
	//   - []commonModules.Project: A slice of Project containing the public projects.
	//   - error: An error object if any error occurs during the database operation.
	GetPublicProjects(limit, offset int) ([]model.Project, error)

	// GetUserProjects retrieves projects associated with a specific user with pagination.
	//
	// It fetches projects created by the user or public projects.
	//
	// Params:
	//   - userID: uint - The ID of the user whose projects are to be retrieved.
	//   - limit: int - The maximum number of projects to retrieve.
	//   - offset: int - The number of projects to skip before starting to collect the result set.
	//
	// Returns:
	//   - []commonModules.Project: A slice of Project containing the user's projects.
	//   - error: An error object if any error occurs during the database operation.
	//
	// Used By:
	// My Projects Page.
	GetUserProjects(userID uint, limit, offset int) ([]model.Project, error)
	// Warning Only For Admin , because it fetches all projects including private ones
	// GetEssentialInfo retrieves essential information of projects with pagination.
	//
	// It selects only the essential fields defined in the Project model.
	//
	// Params:
	//   - limit: int - The maximum number of projects to retrieve.
	//   - offset: int - The number of projects to skip before starting to collect the result set.
	//
	// Returns:
	//   - []dto.ProjectsEssentialInfo: A slice of ProjectsEssentialInfo containing the essential project data.
	//   - error: An error object if any error occurs during the database operation.
	//
	// Used By:
	// Used By Admin.
	GetEssentialInfo(limit, offset int) (*[]model.Project, error)
	GetByID(id uint) (model.Project, error)
	//
	GetProjectsCount() (int, error)
	// IsLiked checks if a project is liked by a user.
	//
	// Params:
	//   - userID: uint - The ID of the user.
	//   - projectID: int - The ID of the project.
	//
	// Returns:
	//   - bool: True if the project is liked by the user, false otherwise.
	//   - error: An error object if any error occurs during the database operation.
	IsLiked(userID uint, projectID uint) (bool, error)
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
	// Create creates a new project in the database.
	//
	// Params:
	//   - project: *commonModules.Project - A pointer to the Project object to be created.
	//
	// Returns:
	//   - error: An error object if any error occurs during the database operation.
	Create(project *model.Project) error
	// LikeProject adds a like from a user to a project.

	LikeProject(userID uint, projectID uint) error
	UnlikeProject(userID uint, projectID uint) error
}

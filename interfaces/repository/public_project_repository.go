package repository

import (
	model "github.com/aruncs31s/esdcmodels"
)

type PublicProjectRepository interface {
	// GetAllProjects retrieves all public projects.
	//
	// Params:
	//  - limit: int - The maximum number of projects to retrieve.
	//  - offset: int - The number of projects to skip before starting to collect the result set.
	// - Returns:
	// - []dto.ProjectResponse - A slice of ProjectResponse DTOs representing the public projects.
	// - error - An error if the retrieval fails.
	GetAllProjects(limit, offset int) (*[]model.Project, error)

	GetUserProjects(user uint, limit, offset int) (*[]model.Project, error)
	// CreateProject is used to create a new project with the provided details.
	//
	// Currently Used In create Project Modal // Remove this comment later
	// Params:
	//  - user: string - The username of the creator.
	// - project: dto.ProjectCreation - The details of the project to be created.
	// - Returns:
	// - *commonModules.Project - The created project.
	// - error - An error if the creation fails.
	GetProject(id uint) (*model.Project, error)
}

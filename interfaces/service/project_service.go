package service

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
)

// TODO: Separate concerns.
type ProjectService interface {
	CreateProject(user string, project dto.ProjectCreation) (*commonModules.Project, error)
	GetProject(id uint, user string) (*dto.ProjectResponse, error)
	// Requested By user for user's own projects
	//
	// Requires authentication
	GetUserProjects(limit, offset int, username string) ([]*dto.ProjectResponse, error)
	ToggleLikeProject(username string, projectID uint) (bool, error) // Returns true if liked, false if unliked
	// UpdateProject(id int, project model.Project) (model.Project, error)
	// DeleteProject(id int) error
}

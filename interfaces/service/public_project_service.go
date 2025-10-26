package service

import "github.com/aruncs31s/esdcprojectmodule/dto"

type PublicProjectService interface {
	GetAllPublicProjects(limit, offset int) (*[]dto.ProjectResponseForPublic, error)
	GetAllUserProjects(username string, limit, offset int) (*[]dto.ProjectResponseForPublic, error)
	GetProject(projectID uint) (*dto.ProjectResponseForPublic, error)
}

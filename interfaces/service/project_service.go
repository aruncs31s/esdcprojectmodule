package service

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
)

type ProjectService interface {
	CreateProject(user string, project dto.ProjectCreation) (*commonModules.Project, error)
	GetProject(id uint, user string) (*dto.ProjectResponse, error)
	GetUserProjects(limit, offset int, username string) ([]*dto.ProjectResponse, error)
	ToggleLikeProject(username string, projectID uint) (bool, error)
	FilterProjects(filter dto.ProjectFilter) (*dto.ProjectListResponse, error)
	UpdateProject(projectID uint, username string, updates dto.ProjectUpdate) error
	DeleteProject(projectID uint, username string, isAdmin bool) error
	GetProjectStats(projectID uint) (*dto.ProjectStats, error)
	CreateComment(username string, comment dto.CommentCreate) (*dto.CommentResponse, error)
	GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error)
	CreateReview(username string, review dto.ReviewCreate) (*dto.ReviewResponse, error)
	GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error)
	DeleteComment(commentID uint, username string, isAdmin bool) error
	ModerateComment(commentID uint, status string) error
}

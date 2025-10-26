package service

import (
	"fmt"

	model "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
	repository "github.com/aruncs31s/esdcprojectmodule/interfaces/repository"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/service"
	utils "github.com/aruncs31s/esdcprojectmodule/utils"
	userRepo "github.com/aruncs31s/esdcusermodule/repository"
)

type publicProjectsService struct {
	publicProjectRepository repository.PublicProjectRepository
	userRepo                userRepo.UserRepository
}

func NewPublicProjectsService(
	publicProjectRepository repository.PublicProjectRepository,
	userRepo userRepo.UserRepository,
) service.PublicProjectService {
	return &publicProjectsService{
		publicProjectRepository: publicProjectRepository,
		userRepo:                userRepo,
	}
}

func (s *publicProjectsService) GetAllPublicProjects(limit, offset int) (*[]dto.ProjectResponseForPublic, error) {
	projects, err := s.publicProjectRepository.GetAllProjects(limit, offset)
	if err != nil {
		return nil, err
	}
	if projects == nil || len(*projects) == 0 {
		return &[]dto.ProjectResponseForPublic{}, nil
	}

	projectsPresentation := getFormatedProjects(projects)
	return &projectsPresentation, nil
}

func (s *publicProjectsService) GetAllUserProjects(user string, limit, offset int) (*[]dto.ProjectResponseForPublic, error) {
	if user == "" {
		return nil, fmt.Errorf("no user specified")
	}

	userID, err := s.userRepo.FindUserIDByUsername(user)
	if err != nil {
		return nil, fmt.Errorf("error fetching user details: %w", err)
	}
	if userID == 0 {
		return nil, fmt.Errorf("invalid user")
	}

	projects, err := s.publicProjectRepository.GetUserProjects(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	projectsPresentation := getFormatedProjects(projects)
	return &projectsPresentation, nil
}
func (s *publicProjectsService) GetProject(projectID uint) (*dto.ProjectResponseForPublic, error) {
	project, err := s.publicProjectRepository.GetProject(projectID)
	if err != nil {
		return nil, err
	}

	projectPresentation := formatProject(project)
	return projectPresentation, nil
}

func getFormatedProjects(projects *[]model.Project) []dto.ProjectResponseForPublic {
	projectsPresentation := make([]dto.ProjectResponseForPublic, len(*projects))
	for i, project := range *projects {
		projectsPresentation[i] = *formatProject(&project)
	}
	return projectsPresentation
}
func formatProject(project *model.Project) *dto.ProjectResponseForPublic {
	return &dto.ProjectResponseForPublic{
		ID:                  project.ID,
		Title:               project.Title,
		Description:         project.Description,
		GithubLink:          project.GithubLink,
		Image:               project.Image,
		LiveUrl:             project.LiveURL,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
		Status:              project.Status,
		Likes:               project.Likes,
		ViewCount:           project.Views,
		CommentCount:        0,
		ForkCount:           0,
		StartCount:          0,
		FavoriteCount:       0,
		Version:             project.Version,
		Cost:                project.Cost,
		Category:            project.Category,
		CreatorDetails:      utils.GetCreatorDetails(project.Creator),
		ContributorsDetails: utils.GetContributorsUsernames(project.Contributors),
		TagsDetails:         utils.GetTagsNames(project.Tags),
		TechnologyDetails:   utils.GetTechnologiesNames(project.Technologies),
	}
}

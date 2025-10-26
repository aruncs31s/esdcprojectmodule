package service

import (
	"fmt"
	"strings"

	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/repository"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/service"
	utils "github.com/aruncs31s/esdcprojectmodule/utils"
	userRepo "github.com/aruncs31s/esdcusermodule/repository"
)

type projectService struct {
	projectRepo repository.ProjectRepository
	userRepo    userRepo.UserRepository
}

func NewProjectService(
	projectRepo repository.ProjectRepository,
	userRepo userRepo.UserRepository,
) service.ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

func (s *projectService) CreateProject(user string, project dto.ProjectCreation) (*commonModules.Project, error) {
	userID, err := s.userRepo.FindUserIDByUsername(user)
	if err != nil {
		return nil, err
	}
	// Build []model.User slice for contributors
	contributors, err := getContributors(s, userID, project)
	if err != nil {
		return nil, err
	}

	// This one , should if the tag exists in the db , if exists assign its value to the project or else create a new tag and assign it to the project
	tags, err := getTags(project, s)
	if err != nil {
		return nil, err
	}
	technologies, err := getTechnologies(project, s)
	if err != nil {
		return nil, err
	}
	// Create the new project
	newProject := commonModules.Project{
		Title:        project.Title,
		Image:        project.Image,
		Description:  project.Description,
		GithubLink:   project.GithubLink,
		Tags:         &tags,
		CreatedBy:    userID,
		ModifiedBy:   &userID,
		Status:       "active", // Default status
		Likes:        0,        // Default value
		Views:        0,
		Category:     project.Category, // Set category from request
		LiveURL:      project.LiveURL,
		Technologies: &technologies,
		Contributors: &contributors,
	}

	// Save the project and its relationships
	if err := s.projectRepo.Create(&newProject); err != nil {
		return nil, err
	}

	return &newProject, nil
}

func getTechnologies(project dto.ProjectCreation, s *projectService) ([]commonModules.Technologies, error) {
	technologies := make([]commonModules.Technologies, 0)
	if project.Technologies != nil {
		splittedTechnologies := strings.Split((*project.Technologies)[0], ",")
		for _, techName := range splittedTechnologies {
			filteredTechName := strings.TrimSpace(techName)
			tech, err := s.projectRepo.FindOrCreateTechnology(filteredTechName)
			if err != nil {
				return nil, fmt.Errorf("error creating/finding technology: %w", err)
			}
			technologies = append(technologies, *tech)
		}
	}
	return technologies, nil
}

func getTags(project dto.ProjectCreation, s *projectService) ([]commonModules.Tag, error) {
	tags := make([]commonModules.Tag, 0)
	if project.Tags != nil {
		for _, tagName := range *project.Tags {
			tag, err := s.projectRepo.FindOrCreateTag(tagName)
			// Check if this error should be avoided.
			if err != nil {
				return nil, fmt.Errorf("error creating/finding tag: %w", err)
			}
			tags = append(tags, *tag)
		}
	}
	return tags, nil
}

func getContributors(s *projectService, userID uint, project dto.ProjectCreation) ([]commonModules.User, error) {
	contributors := make([]commonModules.User, 0)
	// First add the creator as a contributor
	creator, err := s.userRepo.FindByID(userID)
	if err != nil {
		// Move to error class.
		return nil, fmt.Errorf("error fetching creator details: %w", err)
	}

	contributors = append(contributors, *creator)
	// now add contributors from the request
	if project.Contributors != nil && len(*project.Contributors) > 0 {
		users, err := s.userRepo.FindUsersByUsernames(*project.Contributors)
		if err != nil {
			return nil, fmt.Errorf("error fetching contributors: %w", err)
		}
		if len(*users) != len(*project.Contributors) {
			return nil, fmt.Errorf("one or more contributors not found")
		}
		contributors = append(contributors, *users...)
	}
	return contributors, nil
}

func (s *projectService) GetProject(id uint, user string) (*dto.ProjectResponse, error) {
	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	isLiked := false
	if user != "" {
		userID, err := s.userRepo.FindUserIDByUsername(user)
		if err == nil {
			isLiked, _ = s.projectRepo.IsLiked(uint(userID), id)
		}
	}

	p := getProjectResponseForPersonal(project, isLiked)
	return p, nil
}

func getProjectResponseForPersonal(project commonModules.Project, isLiked bool) *dto.ProjectResponse {
	return &dto.ProjectResponse{
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
		Cost:                project.Cost,
		Category:            project.Category,
		IsLiked:             isLiked,
		CreatorDetails:      utils.GetCreatorDetails(project.Creator),
		ContributorsDetails: utils.GetContributorsUsernames(project.Contributors),
		TagsDetails:         utils.GetTagsNames(project.Tags),
		TechnologyDetails:   utils.GetTechnologiesNames(project.Technologies),
	}
}

// HACK: AI
func (s *projectService) ToggleLikeProject(username string, projectID uint) (bool, error) {
	// Get the user ID from username
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return false, err
	}
	// check if already liked.
	isLiked, err := s.projectRepo.IsLiked(userID, projectID)
	if err != nil {
		return false, err
	}
	if isLiked {
		// Unlike
		err = s.projectRepo.UnlikeProject(userID, projectID)
		return false, err
	} else {
		// Like
		err = s.projectRepo.LikeProject(userID, projectID)
		return true, err
	}
}

func (s *projectService) GetUserProjects(limit, offset int, username string) ([]*dto.ProjectResponse, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}
	projects, err := s.projectRepo.GetUserProjects(uint(userID), limit, offset)
	if err != nil {
		return nil, err
	}
	var projectResponses []*dto.ProjectResponse
	for _, project := range projects {
		isLiked, err := s.projectRepo.IsLiked(uint(userID), project.ID)
		if err != nil {
			return nil, err
		}
		p := getProjectResponseForPersonal(project, isLiked)
		projectResponses = append(projectResponses, p)
	}
	return projectResponses, nil
}

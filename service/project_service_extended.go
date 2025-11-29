package service

import (
	"fmt"
	"math"

	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/utils"
)

func (s *projectService) FilterProjects(filter dto.ProjectFilter) (*dto.ProjectListResponse, error) {
	projects, total, err := s.projectRepo.FilterProjects(
		filter.Category,
		filter.Status,
		filter.Visibility,
		filter.Search,
		filter.Technologies,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, err
	}

	projectResponses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		projectResponses[i] = dto.ProjectResponse{
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
			CreatorDetails:      utils.GetCreatorDetails(project.Creator),
			ContributorsDetails: utils.GetContributorsUsernames(project.Contributors),
			TagsDetails:         utils.GetTagsNames(project.Tags),
			TechnologyDetails:   utils.GetTechnologiesNames(project.Technologies),
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	currentPage := (filter.Offset / filter.Limit) + 1

	return &dto.ProjectListResponse{
		Projects: projectResponses,
		Pagination: dto.PaginationMetadata{
			Total:       total,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			PerPage:     filter.Limit,
		},
	}, nil
}

func (s *projectService) UpdateProject(projectID uint, username string, updates dto.ProjectUpdate) error {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return err
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	if project.CreatedBy != userID {
		return fmt.Errorf("unauthorized: only project owner can update")
	}

	updateMap := make(map[string]interface{})
	if updates.Title != nil {
		updateMap["title"] = *updates.Title
	}
	if updates.Description != nil {
		updateMap["description"] = *updates.Description
	}
	if updates.Image != nil {
		updateMap["image"] = *updates.Image
	}
	if updates.Status != nil {
		updateMap["status"] = *updates.Status
	}
	if updates.Visibility != nil {
		updateMap["visibility"] = *updates.Visibility
	}
	if updates.GithubLink != nil {
		updateMap["github_link"] = *updates.GithubLink
	}
	if updates.LiveURL != nil {
		updateMap["live_url"] = *updates.LiveURL
	}
	if updates.Category != nil {
		updateMap["category"] = *updates.Category
	}

	return s.projectRepo.Update(projectID, updateMap)
}

func (s *projectService) DeleteProject(projectID uint, username string, isAdmin bool) error {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return err
	}

	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}

	if !isAdmin && project.CreatedBy != userID {
		return fmt.Errorf("unauthorized: only project owner or admin can delete")
	}

	return s.projectRepo.Delete(projectID, !isAdmin)
}

func (s *projectService) GetProjectStats(projectID uint) (*dto.ProjectStats, error) {
	stats, err := s.projectRepo.GetProjectStats(projectID)
	if err != nil {
		return nil, err
	}

	return &dto.ProjectStats{
		ViewCount:     stats.ViewCount,
		LikeCount:     stats.LikeCount,
		CommentCount:  stats.CommentCount,
		AverageRating: stats.AverageRating,
		ReviewCount:   stats.ReviewCount,
	}, nil
}

func (s *projectService) CreateComment(username string, comment dto.CommentCreate) (*dto.CommentResponse, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	response, err := s.projectRepo.CreateComment(comment.ProjectID, userID, comment.Content)
	if err != nil {
		return nil, err
	}

	user, _ := s.userRepo.FindByID(userID)
	userImage := ""
	if user != nil && user.Image != nil {
		userImage = *user.Image
	}

	return &dto.CommentResponse{
		ID:        response.ID,
		Content:   response.Content,
		ProjectID: response.ProjectID,
		User: dto.Contributor{
			ID:    int(userID),
			Name:  user.Name,
			Email: user.Email,
			Image: userImage,
		},
		CreatedAt: response.CreatedAt,
	}, nil
}

func (s *projectService) GetComments(projectID uint, limit, offset int) ([]dto.CommentResponse, error) {
	comments, err := s.projectRepo.GetComments(projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *projectService) CreateReview(username string, review dto.ReviewCreate) (*dto.ReviewResponse, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	response, err := s.projectRepo.CreateReview(review.ProjectID, userID, review.Rating, review.Comment)
	if err != nil {
		return nil, err
	}

	user, _ := s.userRepo.FindByID(userID)
	userImage := ""
	if user != nil && user.Image != nil {
		userImage = *user.Image
	}

	return &dto.ReviewResponse{
		ID:        response.ID,
		Rating:    response.Rating,
		Comment:   response.Comment,
		ProjectID: response.ProjectID,
		User: dto.Contributor{
			ID:    int(userID),
			Name:  user.Name,
			Email: user.Email,
			Image: userImage,
		},
		CreatedAt: response.CreatedAt,
	}, nil
}

func (s *projectService) GetReviews(projectID uint, limit, offset int) ([]dto.ReviewResponse, error) {
	reviews, err := s.projectRepo.GetReviews(projectID, limit, offset)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (s *projectService) DeleteComment(commentID uint, username string, isAdmin bool) error {
	_, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return err
	}

	// Note: ownership check would require fetching the comment first
	// For now, just allow deletion (or you can add ownership check later)
	return s.projectRepo.DeleteComment(commentID)
}

func (s *projectService) ModerateComment(commentID uint, status string) error {
	return s.projectRepo.UpdateCommentStatus(commentID, status)
}

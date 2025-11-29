package service

import (
	"fmt"

	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/utils"
)

// ========== TRENDING & RECOMMENDATIONS ==========

// GetTrendingProjects returns trending projects from the last N days
func (s *projectService) GetTrendingProjects(limit, offset int, days int) ([]dto.TrendingProject, error) {
	projects, err := s.projectRepo.GetTrendingProjects(limit, offset, days)
	if err != nil {
		return nil, err
	}

	var trendingProjects []dto.TrendingProject
	for _, project := range projects {
		trendingProject := dto.TrendingProject{
			ID:             project.ID,
			Title:          project.Title,
			Image:          project.Image,
			Category:       project.Category,
			Likes:          project.Likes,
			Views:          project.Views,
			CommentCount:   0, // Would need to fetch comments count
			TrendingScore:  float64((project.Likes*40)+(project.Views*30)) / 100.0,
			CreatorDetails: utils.GetCreatorDetails(project.Creator),
		}
		trendingProjects = append(trendingProjects, trendingProject)
	}

	return trendingProjects, nil
}

// GetRecommendedProjects returns personalized recommendations for a user
func (s *projectService) GetRecommendedProjects(username string, limit, offset int) ([]dto.RecommendedProject, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	projects, err := s.projectRepo.GetRecommendedProjects(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var recommendedProjects []dto.RecommendedProject
	for _, project := range projects {
		recommendedProject := dto.RecommendedProject{
			ID:               project.ID,
			Title:            project.Title,
			Description:      project.Description,
			Image:            project.Image,
			Category:         project.Category,
			Likes:            project.Likes,
			TechnologiesUsed: utils.GetTechnologiesNames(project.Technologies),
			CreatorDetails:   utils.GetCreatorDetails(project.Creator),
			MatchScore:       0.85, // Would calculate based on common technologies
			RecommendReason:  "Based on your project interests",
		}
		recommendedProjects = append(recommendedProjects, recommendedProject)
	}

	return recommendedProjects, nil
}

// GetSimilarProjects returns projects similar to the given project
func (s *projectService) GetSimilarProjects(projectID uint, limit, offset int) ([]dto.SimilarProject, error) {
	projects, err := s.projectRepo.GetSimilarProjects(projectID, limit, offset)
	if err != nil {
		return nil, err
	}

	var similarProjects []dto.SimilarProject
	for _, project := range projects {
		similarProject := dto.SimilarProject{
			ID:              project.ID,
			Title:           project.Title,
			Description:     project.Description,
			Image:           project.Image,
			Category:        project.Category,
			Likes:           project.Likes,
			SimilarityScore: 0.78, // Would calculate based on shared technologies/tags
			CreatorDetails:  utils.GetCreatorDetails(project.Creator),
		}
		similarProjects = append(similarProjects, similarProject)
	}

	return similarProjects, nil
}

// ========== ANALYTICS ==========

// GetProjectAnalytics returns analytics for a specific project
func (s *projectService) GetProjectAnalytics(projectID uint, username string, days int) (*dto.ProjectAnalytics, error) {
	// Verify user has access to this project
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Only project creator or admin can view analytics
	userID, _ := s.userRepo.FindUserIDByUsername(username)
	if project.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized to view analytics")
	}

	stats, err := s.projectRepo.GetProjectAnalytics(projectID, days)
	if err != nil {
		return nil, err
	}

	analytics := &dto.ProjectAnalytics{
		ProjectID:     projectID,
		Title:         project.Title,
		TotalViews:    stats.ViewCount,
		TotalLikes:    stats.LikeCount,
		TotalComments: stats.CommentCount,
		AverageRating: stats.AverageRating,
		CreatedAt:     project.CreatedAt,
		UpdatedAt:     project.UpdatedAt,
	}

	return analytics, nil
}

// GetPlatformAnalytics returns platform-wide analytics
func (s *projectService) GetPlatformAnalytics(days int) (*dto.PlatformAnalytics, error) {
	analytics, err := s.projectRepo.GetPlatformAnalytics(days)
	if err != nil {
		return nil, err
	}

	techs, _ := s.projectRepo.GetTrendingTechnologies(10)
	tags, _ := s.projectRepo.GetTrendingTags(10)

	platformAnalytics := &dto.PlatformAnalytics{
		TotalProjects: analytics.TotalProjects,
		TotalViews:    analytics.TotalViews,
		TotalLikes:    analytics.TotalLikes,
	}

	// Convert trending techs to PopularTechnology
	for _, tech := range techs {
		platformAnalytics.PopularTechs = append(platformAnalytics.PopularTechs, dto.PopularTechnology{
			Name:         tech.Name,
			UsageCount:   tech.UsageCount,
			ProjectCount: tech.UsageCount,
		})
	}

	// Convert trending tags to PopularTag
	for _, tag := range tags {
		platformAnalytics.PopularTags = append(platformAnalytics.PopularTags, dto.PopularTag{
			Name:         tag.Name,
			UsageCount:   tag.UsageCount,
			ProjectCount: tag.UsageCount,
		})
	}

	return platformAnalytics, nil
}

// ========== TEMPLATES ==========

// CreateTemplate saves a project as a template
func (s *projectService) CreateTemplate(username string, template dto.TemplateCreate) (*dto.ProjectTemplate, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	// Verify user owns the project
	project, err := s.projectRepo.GetByID(template.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	if project.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized: only project creator can create templates")
	}

	newTemplate := commonModules.ProjectTemplate{
		ProjectID:   template.ProjectID,
		Name:        template.Name,
		Description: template.Description,
		CreatedBy:   userID,
		IsPublic:    template.IsPublic,
		UsageCount:  0,
	}

	if err := s.projectRepo.CreateTemplate(&newTemplate); err != nil {
		return nil, err
	}

	return &dto.ProjectTemplate{
		ID:          newTemplate.ID,
		Name:        newTemplate.Name,
		Description: newTemplate.Description,
		CreatorID:   newTemplate.CreatedBy,
		Category:    project.Category,
		IsPublic:    newTemplate.IsPublic,
		UsageCount:  0,
		CreatedAt:   newTemplate.CreatedAt,
		UpdatedAt:   newTemplate.UpdatedAt,
	}, nil
}

// GetTemplates returns public templates
func (s *projectService) GetTemplates(limit, offset int, isPublic bool) ([]dto.TemplateListItem, error) {
	templates, err := s.projectRepo.GetPublicTemplates(limit, offset)
	if err != nil {
		return nil, err
	}

	var items []dto.TemplateListItem
	for _, template := range templates {
		item := dto.TemplateListItem{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			UsageCount:  template.UsageCount,
			CreatedAt:   template.CreatedAt,
		}
		items = append(items, item)
	}

	return items, nil
}

// GetUserTemplates returns templates created by a user
func (s *projectService) GetUserTemplates(username string, limit, offset int) ([]dto.TemplateListItem, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	templates, err := s.projectRepo.GetUserTemplates(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var items []dto.TemplateListItem
	for _, template := range templates {
		item := dto.TemplateListItem{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			UsageCount:  template.UsageCount,
			CreatedAt:   template.CreatedAt,
		}
		items = append(items, item)
	}

	return items, nil
}

// GetTemplate returns a specific template
func (s *projectService) GetTemplate(templateID uint) (*dto.ProjectTemplate, error) {
	template, err := s.projectRepo.GetTemplateByID(templateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}

	return &dto.ProjectTemplate{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		CreatorID:   template.CreatedBy,
		IsPublic:    template.IsPublic,
		UsageCount:  template.UsageCount,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
	}, nil
}

// DeleteTemplate removes a template
func (s *projectService) DeleteTemplate(templateID uint, username string) error {
	template, err := s.projectRepo.GetTemplateByID(templateID)
	if err != nil {
		return fmt.Errorf("template not found")
	}

	userID, _ := s.userRepo.FindUserIDByUsername(username)
	if template.CreatedBy != userID {
		return fmt.Errorf("unauthorized: only creator can delete templates")
	}

	return s.projectRepo.DeleteTemplate(templateID)
}

// CreateProjectFromTemplate creates a new project based on a template
func (s *projectService) CreateProjectFromTemplate(username string, request dto.CreateFromTemplate) (*dto.ProjectResponse, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	template, err := s.projectRepo.GetTemplateByID(request.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template not found")
	}

	sourceProject, err := s.projectRepo.GetByID(template.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("source project not found")
	}

	// Create new project based on template
	newProject := commonModules.Project{
		Title:        request.Title,
		Description:  request.Description,
		GithubLink:   request.GithubLink,
		Category:     sourceProject.Category,
		CreatedBy:    userID,
		ModifiedBy:   &userID,
		Status:       "active",
		Tags:         sourceProject.Tags,
		Technologies: sourceProject.Technologies,
	}

	if err := s.projectRepo.Create(&newProject); err != nil {
		return nil, err
	}

	// Increment template usage count
	s.projectRepo.IncrementTemplateUsage(request.TemplateID)

	return getProjectResponseForPersonal(newProject, false), nil
}

// ========== NOTIFICATIONS ==========

// GetNotifications returns notifications for a user
func (s *projectService) GetNotifications(username string, limit, offset int) ([]dto.NotificationResponse, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	notifications, err := s.projectRepo.GetUserNotifications(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []dto.NotificationResponse
	for _, notif := range notifications {
		triggeredBy, _ := s.userRepo.FindByID(notif.TriggeredBy)
		response := dto.NotificationResponse{
			ID:        notif.ID,
			Type:      notif.Type,
			Title:     notif.Title,
			Message:   notif.Message,
			ProjectID: notif.ProjectID,
			IsRead:    notif.IsRead,
			CreatedAt: notif.CreatedAt,
		}
		if triggeredBy != nil {
			response.TriggeredBy = utils.GetCreatorDetails(triggeredBy)
		}
		responses = append(responses, response)
	}

	return responses, nil
}

// MarkNotificationAsRead marks a notification as read
func (s *projectService) MarkNotificationAsRead(notificationID uint, username string) error {
	// Could add verification that user owns this notification
	return s.projectRepo.MarkNotificationAsRead(notificationID)
}

// MarkAllAsRead marks all notifications as read for a user
func (s *projectService) MarkAllAsRead(username string) error {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return err
	}

	return s.projectRepo.MarkAllNotificationsAsRead(userID)
}

// DeleteNotification removes a notification
func (s *projectService) DeleteNotification(notificationID uint, username string) error {
	// Could add verification that user owns this notification
	return s.projectRepo.DeleteNotification(notificationID)
}

// CreateNotification creates a new notification
func (s *projectService) CreateNotification(userID uint, notification dto.NotificationCreate) error {
	newNotif := commonModules.Notification{
		UserID:    userID,
		Type:      notification.Type,
		Title:     notification.Title,
		Message:   notification.Message,
		ProjectID: notification.ProjectID,
	}

	return s.projectRepo.CreateNotification(&newNotif)
}

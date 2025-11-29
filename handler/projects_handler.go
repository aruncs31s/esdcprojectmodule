package handler

import (
	"strings"

	"github.com/aruncs31s/esdcsharedhelpersmodule/helper"
	sharedHelper "github.com/aruncs31s/esdcsharedhelpersmodule/interface/helper"

	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/service"
	"github.com/aruncs31s/responsehelper"
	"github.com/gin-gonic/gin"
)

type ProjectHandler interface {
	// CreateProject is used to create a new project.
	//
	// Used by Private Routes only.
	// Requires authentication.
	CreateProject(c *gin.Context)
	GetAllProjects(c *gin.Context)
	// GetProject is used to retrieve a project by its ID.
	GetProject(c *gin.Context)
	// ToggleLikeProject is used to like or unlike a project.
	ToggleLikeProject(c *gin.Context)
	// Trending & Recommendations
	GetTrendingProjects(c *gin.Context)
	GetRecommendedProjects(c *gin.Context)
	GetSimilarProjects(c *gin.Context)
	// Analytics
	GetProjectAnalytics(c *gin.Context)
	GetPlatformAnalytics(c *gin.Context)
	// Templates
	CreateTemplate(c *gin.Context)
	GetTemplates(c *gin.Context)
	GetUserTemplates(c *gin.Context)
	GetTemplate(c *gin.Context)
	DeleteTemplate(c *gin.Context)
	CreateProjectFromTemplate(c *gin.Context)
	// Notifications
	GetNotifications(c *gin.Context)
	MarkNotificationAsRead(c *gin.Context)
	MarkAllNotificationsAsRead(c *gin.Context)
	DeleteNotification(c *gin.Context)
	// Export
	ExportProject(c *gin.Context)
	ExportPortfolio(c *gin.Context)
	// UpdateProject is used to update an existing project.
	// UpdateProject(c *gin.Context)
	// // DeleteProject is used to delete a project.
	// DeleteProject(c *gin.Context)
}

type projectHandler struct {
	responseHelper responsehelper.ResponseHelper
	projectService service.ProjectService
	requestHelper  sharedHelper.RequestHelper
	validator      sharedHelper.RequestValidator
}

func NewProjectHandler(projectService service.ProjectService) ProjectHandler {
	responseHelper, requestHelper, validator := getHelpers()
	return &projectHandler{
		responseHelper: responseHelper,
		projectService: projectService,
		requestHelper:  requestHelper,
		validator:      validator,
	}
}

func (h *projectHandler) GetValidator() sharedHelper.RequestValidator {
	return h.validator
}
func (h *projectHandler) GetResponseHelper() responsehelper.ResponseHelper {
	return h.responseHelper
}

// CreateProject Is used to create a new project.
// Used by Private Routes only.
// Requires authentication.
func (h *projectHandler) CreateProject(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}
	projectData, failed := helper.GetJSONDataFromRequest[dto.ProjectCreation](c, h.responseHelper)
	if failed {
		return
	}
	createdProject, err := h.projectService.CreateProject(user, projectData)
	// move these errors to a common class
	if err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed") {
		h.responseHelper.BadRequest(c, "Project with the same name already exists", err.Error())
		return
	}
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to create project", err)
		return
	}
	h.responseHelper.Created(c, createdProject)
}
func (h *projectHandler) GetAllProjects(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}
	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	projects, err := h.projectService.GetUserProjects(limit, offset, user)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to retrieve projects", err)
		return
	}
	h.responseHelper.Success(c, projects)
}
func (h *projectHandler) GetProject(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}
	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide an id.")
	if failed {
		return
	}

	project, err := h.projectService.GetProject(id, user)
	if err != nil {
		h.responseHelper.NotFound(c, "Project not found")
		return
	}
	h.responseHelper.Success(c, project)
}

// ToggleLikeProject godoc
// @Summary Toggle like on a project
// @Description Like or unlike a project
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]interface{} "Like toggled successfully"
// @Failure 400 {object} map[string]interface{} "Invalid project ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Project not found"
// @Router /projects/{id}/like [post]
func (h *projectHandler) ToggleLikeProject(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}
	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide an id.")
	if failed {
		return
	}
	liked, err := h.projectService.ToggleLikeProject(user, id)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to toggle like", err)
		return
	}
	response := map[string]interface{}{
		"liked":   liked,
		"message": "Like toggled successfully",
	}
	h.responseHelper.Success(c, response)
}

// ========== TRENDING & RECOMMENDATIONS ==========

// GetTrendingProjects returns trending projects
func (h *projectHandler) GetTrendingProjects(c *gin.Context) {
	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	days := 30 // Default to last 30 days

	projects, err := h.projectService.GetTrendingProjects(limit, offset, days)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch trending projects", err)
		return
	}
	h.responseHelper.Success(c, projects)
}

// GetRecommendedProjects returns personalized recommendations
func (h *projectHandler) GetRecommendedProjects(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	projects, err := h.projectService.GetRecommendedProjects(user, limit, offset)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch recommendations", err)
		return
	}
	h.responseHelper.Success(c, projects)
}

// GetSimilarProjects returns projects similar to the specified project
func (h *projectHandler) GetSimilarProjects(c *gin.Context) {
	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide an id.")
	if failed {
		return
	}

	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	projects, err := h.projectService.GetSimilarProjects(id, limit, offset)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch similar projects", err)
		return
	}
	h.responseHelper.Success(c, projects)
}

// ========== ANALYTICS ==========

// GetProjectAnalytics returns analytics for a project
func (h *projectHandler) GetProjectAnalytics(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide an id.")
	if failed {
		return
	}

	analytics, err := h.projectService.GetProjectAnalytics(id, user, 30)
	if err != nil {
		h.responseHelper.BadRequest(c, "Failed to fetch analytics", err.Error())
		return
	}
	h.responseHelper.Success(c, analytics)
}

// GetPlatformAnalytics returns platform-wide analytics
func (h *projectHandler) GetPlatformAnalytics(c *gin.Context) {
	analytics, err := h.projectService.GetPlatformAnalytics(30)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch platform analytics", err)
		return
	}
	h.responseHelper.Success(c, analytics)
}

// ========== TEMPLATES ==========

// CreateTemplate creates a project template
func (h *projectHandler) CreateTemplate(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	var req dto.TemplateCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseHelper.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	template, err := h.projectService.CreateTemplate(user, req)
	if err != nil {
		h.responseHelper.BadRequest(c, "Failed to create template", err.Error())
		return
	}
	h.responseHelper.Created(c, template)
}

// GetTemplates returns public templates
func (h *projectHandler) GetTemplates(c *gin.Context) {
	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	templates, err := h.projectService.GetTemplates(limit, offset, true)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch templates", err)
		return
	}
	h.responseHelper.Success(c, templates)
}

// GetUserTemplates returns templates created by the current user
func (h *projectHandler) GetUserTemplates(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	templates, err := h.projectService.GetUserTemplates(user, limit, offset)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch templates", err)
		return
	}
	h.responseHelper.Success(c, templates)
}

// GetTemplate returns a specific template
func (h *projectHandler) GetTemplate(c *gin.Context) {
	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide a template id.")
	if failed {
		return
	}

	template, err := h.projectService.GetTemplate(id)
	if err != nil {
		h.responseHelper.NotFound(c, "Template not found")
		return
	}
	h.responseHelper.Success(c, template)
}

// DeleteTemplate deletes a template
func (h *projectHandler) DeleteTemplate(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide a template id.")
	if failed {
		return
	}

	err := h.projectService.DeleteTemplate(id, user)
	if err != nil {
		h.responseHelper.BadRequest(c, "Failed to delete template", err.Error())
		return
	}
	h.responseHelper.Success(c, map[string]string{"message": "Template deleted successfully"})
}

// CreateProjectFromTemplate creates a project from a template
func (h *projectHandler) CreateProjectFromTemplate(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	var req dto.CreateFromTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		h.responseHelper.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	project, err := h.projectService.CreateProjectFromTemplate(user, req)
	if err != nil {
		h.responseHelper.BadRequest(c, "Failed to create project from template", err.Error())
		return
	}
	h.responseHelper.Created(c, project)
}

// ========== NOTIFICATIONS ==========

// GetNotifications returns user notifications
func (h *projectHandler) GetNotifications(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	notifications, err := h.projectService.GetNotifications(user, limit, offset)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to fetch notifications", err)
		return
	}
	h.responseHelper.Success(c, notifications)
}

// MarkNotificationAsRead marks a notification as read
func (h *projectHandler) MarkNotificationAsRead(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide a notification id.")
	if failed {
		return
	}

	err := h.projectService.MarkNotificationAsRead(id, user)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to mark notification as read", err)
		return
	}
	h.responseHelper.Success(c, map[string]string{"message": "Notification marked as read"})
}

// MarkAllNotificationsAsRead marks all notifications as read
func (h *projectHandler) MarkAllNotificationsAsRead(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	err := h.projectService.MarkAllAsRead(user)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to mark notifications as read", err)
		return
	}
	h.responseHelper.Success(c, map[string]string{"message": "All notifications marked as read"})
}

// DeleteNotification deletes a notification
func (h *projectHandler) DeleteNotification(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	id, failed := h.requestHelper.ValidateAndParseID(h, "id", c, "please provide a notification id.")
	if failed {
		return
	}

	err := h.projectService.DeleteNotification(id, user)
	if err != nil {
		h.responseHelper.InternalError(c, "Failed to delete notification", err)
		return
	}
	h.responseHelper.Success(c, map[string]string{"message": "Notification deleted"})
}

// ========== EXPORT ==========

// ExportProject exports a project in the specified format
func (h *projectHandler) ExportProject(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	_, failed = h.requestHelper.ValidateAndParseID(h, "id", c, "please provide an id.")
	if failed {
		return
	}

	format := c.DefaultQuery("format", "json")

	// Export logic would be implemented here
	h.responseHelper.Success(c, map[string]string{
		"message": "Export functionality requires exportService implementation",
		"format":  format,
		"user":    user,
	})
}

// ExportPortfolio exports user's portfolio
func (h *projectHandler) ExportPortfolio(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}

	format := c.DefaultQuery("format", "json")
	includeStats := c.DefaultQuery("include_stats", "false") == "true"

	// Export logic would be implemented here
	h.responseHelper.Success(c, map[string]interface{}{
		"message":       "Portfolio export functionality requires exportService implementation",
		"format":        format,
		"include_stats": includeStats,
		"username":      user,
	})
}

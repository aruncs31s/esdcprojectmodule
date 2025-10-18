package handler

import (
	"strconv"
	"strings"

	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/service"
	"github.com/aruncs31s/responsehelper"
	"github.com/gin-gonic/gin"
)

type ProjectHandler interface {
	GetPublicProjects(c *gin.Context)
	GetUserProjects(c *gin.Context)
	CreateProject(c *gin.Context)
	GetProject(c *gin.Context)
	ToggleLikeProject(c *gin.Context)
	// UpdateProject(c *gin.Context)
	// DeleteProject(c *gin.Context)
}
type projectHandler struct {
	responseHelper responsehelper.ResponseHelper
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) ProjectHandler {
	responseHelper := responsehelper.NewResponseHelper()
	return &projectHandler{
		responseHelper: responseHelper,
		projectService: projectService,
	}
}

// GetPublicProjects godoc
// @Summary Get public projects
// @Description Retrieve all public projects (no authentication required)
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Public projects retrieved successfully"
// @Failure 404 {object} map[string]interface{} "No public projects found"
// @Router /projects/public [get]
func (h *projectHandler) GetPublicProjects(c *gin.Context) {
	// Pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	user := c.GetString("user") // May be empty if not authenticated

	projects, err := h.projectService.GetPublicProjects(limit, offset, user)
	if err != nil {
		h.responseHelper.NotFound(c, "No public projects found")
		return
	}
	h.responseHelper.Success(c, projects)
}

// GetUserProjects godoc
// @Summary Get user's projects
// @Description Retrieve projects visible to the authenticated user (public + own projects)
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User projects retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "No projects found"
// @Router /projects/user [get]
func (h *projectHandler) GetUserProjects(c *gin.Context) {
	user := c.GetString("user")
	if user == "" {
		h.responseHelper.Unauthorized(c, "Authentication required")
		return
	}

	// Pagination parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	projects, err := h.projectService.GetUserProjects(limit, offset, user)
	if err != nil {
		h.responseHelper.NotFound(c, "No projects found")
		return
	}
	h.responseHelper.Success(c, projects)
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project with the provided data
// @Tags projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project body dto.ProjectCreation true "Project creation data"
// @Success 201 {object} map[string]interface{} "Project created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Failure 500 {object} map[string]interface{} "Failed to create project"
// @Router /projects [post]
func (h *projectHandler) CreateProject(c *gin.Context) {
	user := c.GetString("user")
	if user == "" {
		h.responseHelper.Unauthorized(c, "User not authenticated")
		return
	}
	var project dto.ProjectCreation
	if err := c.ShouldBindJSON(&project); err != nil {
		h.responseHelper.BadRequest(c, "Invalid request body", err.Error())
		return
	}
	createdProject, err := h.projectService.CreateProject(user, project)
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
func (h *projectHandler) GetProject(c *gin.Context) {
	// user := c.GetString("user")
	idStr := c.Param("id")
	if idStr == "" {
		h.responseHelper.BadRequest(c, "Project ID is required", "invalid id")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseHelper.BadRequest(c, "Invalid project ID", err.Error())
		return
	}

	user := c.GetString("user") // May be empty

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
	user := c.GetString("user")
	if user == "" {
		h.responseHelper.Unauthorized(c, "User not authenticated")
		return
	}
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.responseHelper.BadRequest(c, "Invalid project ID", err.Error())
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

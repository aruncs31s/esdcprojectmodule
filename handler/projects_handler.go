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

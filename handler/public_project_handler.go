// Package handler provides HTTP handlers for managing public and user-specific projects.
//
// This package defines the `publicProjectHandler` struct, which implements the
// `PublicProjectHandler` interface. It provides methods for handling requests
// related to public and user-specific projects.
//
// The handler depends on the following interfaces:
// - `service.PublicProjectService`: Provides methods to retrieve public and user-specific projects.
// - `sharedHelper.RequestHelper`: Provides utility methods for extracting request parameters.
// - `responsehelper.ResponseHelper`: Provides methods for sending HTTP responses.
// - `sharedHelper.RequestValidator`: Provides methods for validating incoming requests.
//
// The `NewPublicProjectHandler` function initializes and returns a new instance of
// `publicProjectHandler` with the required dependencies.
//
// The following methods are implemented:
//   - `GetValidator`: Returns the request validator instance.
//   - `GetResponseHelper`: Returns the response helper instance.
//   - `GetPublicProjects`: Handles requests to retrieve public projects. This method
//     does not require authentication.
//   - `GetUserProjects`: Handles requests to retrieve projects visible to the authenticated
//     user (public + own projects). This method requires authentication.
//
// Usage:
//
// The `GetPublicProjects` method is intended for public routes and does not require
// authentication. It retrieves a paginated list of public projects.
//
// The `GetUserProjects` method retrieves a paginated list of projects visible to the
// authenticated user. It requires a valid Bearer token for authentication.
//
// Example:
//
// To use this handler, initialize it with the required dependencies and register
// its methods with the appropriate routes in a Gin router.
package handler

import (
	"github.com/aruncs31s/esdcprojectmodule/interfaces/handler"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/service"
	sharedHelperImpl "github.com/aruncs31s/esdcsharedhelpersmodule/helper"
	sharedHelper "github.com/aruncs31s/esdcsharedhelpersmodule/interface/helper"
	"github.com/aruncs31s/esdcsharedhelpersmodule/utils"
	"github.com/aruncs31s/responsehelper"
	"github.com/gin-gonic/gin"
)

type publicProjectHandler struct {
	publicProjectService service.PublicProjectService
	requestHelper        sharedHelper.RequestHelper
	responseHelper       responsehelper.ResponseHelper
	validator            sharedHelper.RequestValidator
}

func (h *publicProjectHandler) GetValidator() sharedHelper.RequestValidator {
	return h.validator
}
func (h *publicProjectHandler) GetResponseHelper() responsehelper.ResponseHelper {
	return h.responseHelper
}

func NewPublicProjectHandler(publicProjectService service.PublicProjectService) handler.PublicProjectHandler {
	responseHelper, requestHelper, validator := getHelpers()
	return &publicProjectHandler{
		publicProjectService: publicProjectService,
		requestHelper:        requestHelper,
		responseHelper:       responseHelper,
		validator:            validator,
	}
}

func getHelpers() (responsehelper.ResponseHelper, sharedHelper.RequestHelper, sharedHelper.RequestValidator) {
	responseHelper := responsehelper.NewResponseHelper()
	requestHelper := sharedHelperImpl.NewRequestHelper()
	validator := sharedHelperImpl.NewRequestValidator()
	return responseHelper, requestHelper, validator
}

// GetPublicProjects Must Be used by Public Routes only.
// This does not require authentication.
func (h *publicProjectHandler) GetPublicProjects(c *gin.Context) {
	// Pagination parameters
	limit, offset := h.requestHelper.GetLimitAndOffset(c)
	projects, err := h.publicProjectService.GetAllPublicProjects(limit, offset)
	if err != nil {
		h.responseHelper.NotFound(c, "No public projects found")
		return
	}
	h.responseHelper.Success(c, projects)
}

func (h *publicProjectHandler) GetUserProjects(c *gin.Context) {
	user, failed := h.requestHelper.GetAndValidateUsername(c, h)
	if failed {
		return
	}
	// Pagination parameters
	limit, offset := h.requestHelper.GetLimitAndOffset(c)

	projects, err := h.publicProjectService.GetAllUserProjects(user, limit, offset)
	if err != nil {
		h.responseHelper.NotFound(c, "No projects found")
		return
	}
	h.responseHelper.Success(c, projects)
}
func (h *publicProjectHandler) GetProject(c *gin.Context) {
	projectID, failed := h.requestHelper.ValidateAndParseID(h, "id", c, utils.ErrBadRequest.Error())
	if failed {
		return
	}
	project, err := h.publicProjectService.GetProject(projectID)
	if err != nil {
		h.responseHelper.NotFound(c, "Project not found")
		return
	}
	h.responseHelper.Success(c, project)
}

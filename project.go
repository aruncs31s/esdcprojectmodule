package project

import (
	"github.com/aruncs31s/esdcprojectmodule/handler"
	"github.com/aruncs31s/esdcprojectmodule/repository"
	"github.com/aruncs31s/esdcprojectmodule/routes"
	"github.com/aruncs31s/esdcprojectmodule/service"
	userRepo "github.com/aruncs31s/esdcusermodule/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type projectModule struct {
	projectHandler handler.ProjectHandler
	r              *gin.Engine
}

var projectInstance *projectModule

// InitProjectModule initializes the project module with the provided Gin engine and GORM database.
//
// It creates instances of repositories, services, and handlers, and sets up the module instance.
//
// Params:
//   - r: *gin.Engine - The Gin engine to register routes on.
//
// - db: *gorm.DB - The GORM database connection.

func InitProjectModule(r *gin.Engine, db *gorm.DB) {
	projectRepository := repository.NewProjectRepository(db)
	userRepository := userRepo.NewUserRepository(db)
	projectService := service.NewProjectService(projectRepository, userRepository)
	projectHandler := handler.NewProjectHandler(projectService)
	projectInstance = &projectModule{
		projectHandler: projectHandler,
		r:              r,
	}
}

// RegisterPublicProjectRoutes registers the public project routes with the Gin engine.
//
// It sets up the routes that are accessible without authentication.

func RegisterPublicProjectRoutes() {
	routes.RegisterPublicProjectRoutes(projectInstance.r, projectInstance.projectHandler)
}

// RegisterPrivateProjectRoutes registers the private project routes with the Gin engine.
//
// It sets up the routes that require authentication.
//
// Params:
// - r: *gin.Engine - The Gin engine to register routes on.
//
// Note: Only Use this after enabling jwt middleware on the routes.
func RegisterPrivateProjectRoutes(r *gin.Engine) {
	routes.RegisterPrivateProjectRoutes(r, projectInstance.projectHandler)
}

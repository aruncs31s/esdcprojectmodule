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

func RegisterPublicProjectRoutes() {
	routes.RegisterPublicProjectRoutes(projectInstance.r, projectInstance.projectHandler)
}
func RegisterPrivateProjectRoutes(r *gin.Engine) {
	routes.RegisterPrivateProjectRoutes(r, projectInstance.projectHandler)
}

package handler

import "github.com/gin-gonic/gin"

type PublicProjectHandler interface {
	// GetPublicProjects handles the retrieval of public projects.
	// This method is intended for public routes and does not require authentication.
	// It retrieves a paginated list of public projects and sends them in the response.
	//
	// Parameters:
	// - c: The Gin context, which contains the HTTP request and response.
	//
	// Behavior:
	// - Extracts pagination parameters (limit and offset) from the request.
	// - Calls the service layer to fetch public projects.
	// - If no projects are found, responds with a 404 status and an appropriate message.
	// - If projects are found, responds with a 200 status and the list of projects.
	GetPublicProjects(c *gin.Context)
	// When someone clicks a user , this handler gets called to get the user projects , still no private projects are send.
	GetUserProjects(c *gin.Context)
	// GetProject (by id) for both public projects and private projects
	// In future implement a method method to fetch the private projects also if the request is comming from an authenticated user.
	GetProject(c *gin.Context)
}

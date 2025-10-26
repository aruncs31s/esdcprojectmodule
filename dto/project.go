package dto

import "time"

// ProjectCreation represents project creation request
// @Description Project creation request payload

type ProjectCreation struct {
	Title        string    `json:"title" example:"My Project"`
	Image        *string   `json:"image" example:"https://example.com/image.jpg"`
	Description  string    `json:"description" example:"This is a sample project description"`
	Status       string    `json:"status" example:"in_progress"`
	Visibility   string    `json:"visibility" example:"everyone"`
	GithubLink   string    `json:"github_link" example:"https://github.com/user/project"`
	Technologies *[]string `json:"technologies" example:"Go, Gin, GORM"`
	Tags         *[]string `json:"tags" example:"backend,api"`
	LiveURL      *string   `json:"live_url" example:"https://example.com/live"`
	Category     string    `json:"category" example:"Web Development"`
	Contributors *[]string `json:"contributors" example:"2,3,4"`
}

type ProjectResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       *string   `json:"image"`
	GithubLink  string    `json:"github_link"`
	LiveUrl     *string   `json:"live_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	// Newly Addedd
	Likes               int            `json:"likes"`
	Cost                int            `json:"cost"`
	Category            string         `json:"category"`
	IsLiked             bool           `json:"is_liked"`
	CreatorDetails      Contributor    `json:"creator_details,omitempty"`
	ContributorsDetails *[]Contributor `json:"contributors_details,omitempty"`
	TagsDetails         *[]Tag         `json:"tags_details,omitempty"`
	TechnologyDetails   *[]Technology  `json:"technology_details,omitempty"`
}

type ProjectResponseForPublic struct {
	ID                  uint           `json:"id"`
	Title               string         `json:"title"`
	Description         string         `json:"description"`
	Image               *string        `json:"image"`
	GithubLink          string         `json:"github_link"`
	LiveUrl             *string        `json:"live_url"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	Status              string         `json:"status"`
	Likes               int            `json:"likes"`
	ViewCount           int            `json:"view_count"`
	ForkCount           int            `json:"fork_count"`
	CommentCount        int            `json:"comment_count"`
	StartCount          int            `json:"star_count"`
	FavoriteCount       int            `json:"favorite_count"`
	Version             string         `json:"version"`
	Cost                int            `json:"cost"`
	Category            string         `json:"category"`
	CreatorDetails      Contributor    `json:"creator_details,omitempty"`
	ContributorsDetails *[]Contributor `json:"contributors_details,omitempty"`
	TagsDetails         *[]Tag         `json:"tags_details,omitempty"`
	TechnologyDetails   *[]Technology  `json:"technology_details,omitempty"`
	LikedBy             []User         `json:"liked_by,omitempty"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}

type Contributor struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
}
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Technology struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// For Admin Pannel
type ProjectsEssentialInfo struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	CreatedBy  string `json:"created_by"`
	Image      string `json:"image"`
	Status     string `json:"status"`
	Visibility string `json:"visibility"`
	Likes      int    `json:"likes"`
	Views      int    `json:"views"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

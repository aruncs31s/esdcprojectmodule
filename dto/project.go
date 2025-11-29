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

// ========== TRENDING & RECOMMENDATIONS ==========

type TrendingProject struct {
	ID             uint        `json:"id"`
	Title          string      `json:"title"`
	Image          *string     `json:"image"`
	Category       string      `json:"category"`
	Likes          int         `json:"likes"`
	Views          int         `json:"views"`
	CommentCount   int         `json:"comment_count"`
	TrendingScore  float64     `json:"trending_score"`
	CreatorDetails Contributor `json:"creator_details,omitempty"`
}

type RecommendedProject struct {
	ID               uint          `json:"id"`
	Title            string        `json:"title"`
	Description      string        `json:"description"`
	Image            *string       `json:"image"`
	Category         string        `json:"category"`
	Likes            int           `json:"likes"`
	TechnologiesUsed *[]Technology `json:"technologies_used,omitempty"`
	CreatorDetails   Contributor   `json:"creator_details,omitempty"`
	MatchScore       float64       `json:"match_score"`
	RecommendReason  string        `json:"recommend_reason"`
}

type SimilarProject struct {
	ID              uint        `json:"id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	Image           *string     `json:"image"`
	Category        string      `json:"category"`
	Likes           int         `json:"likes"`
	SimilarityScore float64     `json:"similarity_score"`
	CommonTags      []string    `json:"common_tags"`
	CreatorDetails  Contributor `json:"creator_details,omitempty"`
}

// ========== ANALYTICS ==========

type ProjectAnalytics struct {
	ProjectID           uint          `json:"project_id"`
	Title               string        `json:"title"`
	TotalViews          int           `json:"total_views"`
	TotalLikes          int           `json:"total_likes"`
	TotalComments       int           `json:"total_comments"`
	AverageRating       float64       `json:"average_rating"`
	ViewsTrend          []DailyMetric `json:"views_trend"`
	LikesTrend          []DailyMetric `json:"likes_trend"`
	PopularTechnologies []TechMetric  `json:"popular_technologies"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
}

type DailyMetric struct {
	Date   time.Time `json:"date"`
	Count  int       `json:"count"`
	Change float64   `json:"change_percent"`
}

type TechMetric struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type PopularTechnology struct {
	Name         string `json:"name"`
	UsageCount   int    `json:"usage_count"`
	ProjectCount int    `json:"project_count"`
}

type PopularTag struct {
	Name         string `json:"name"`
	UsageCount   int    `json:"usage_count"`
	ProjectCount int    `json:"project_count"`
}

type PlatformAnalytics struct {
	TotalProjects        int                 `json:"total_projects"`
	TotalViews           int                 `json:"total_views"`
	TotalLikes           int                 `json:"total_likes"`
	PopularTechs         []PopularTechnology `json:"popular_technologies"`
	PopularTags          []PopularTag        `json:"popular_tags"`
	TopProjects          []TrendingProject   `json:"top_projects"`
	CategoryDistribution []CategoryMetric    `json:"category_distribution"`
}

type CategoryMetric struct {
	Category     string  `json:"category"`
	ProjectCount int     `json:"project_count"`
	Percentage   float64 `json:"percentage"`
}

// ========== NOTIFICATIONS ==========

type Notification struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Type        string    `json:"type"` // like, comment, follow, milestone
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	ProjectID   uint      `json:"project_id"`
	TriggeredBy uint      `json:"triggered_by"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type NotificationResponse struct {
	ID          uint        `json:"id"`
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Message     string      `json:"message"`
	ProjectID   uint        `json:"project_id"`
	TriggeredBy Contributor `json:"triggered_by"`
	IsRead      bool        `json:"is_read"`
	CreatedAt   time.Time   `json:"created_at"`
}

type NotificationCreate struct {
	Type      string `json:"type" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Message   string `json:"message" binding:"required"`
	ProjectID uint   `json:"project_id"`
}

// ========== PROJECT TEMPLATES ==========

type ProjectTemplate struct {
	ID             uint          `json:"id"`
	Name           string        `json:"name"`
	Description    string        `json:"description"`
	Image          *string       `json:"image"`
	CreatorID      uint          `json:"creator_id"`
	CreatorDetails Contributor   `json:"creator_details,omitempty"`
	Technologies   *[]Technology `json:"technologies,omitempty"`
	Tags           *[]Tag        `json:"tags,omitempty"`
	Category       string        `json:"category"`
	IsPublic       bool          `json:"is_public"`
	UsageCount     int           `json:"usage_count"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type TemplateCreate struct {
	ProjectID   uint   `json:"project_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	IsPublic    bool   `json:"is_public"`
}

type TemplateListItem struct {
	ID             uint        `json:"id"`
	Name           string      `json:"name"`
	Description    string      `json:"description"`
	Image          *string     `json:"image"`
	Category       string      `json:"category"`
	UsageCount     int         `json:"usage_count"`
	CreatorDetails Contributor `json:"creator_details,omitempty"`
	CreatedAt      time.Time   `json:"created_at"`
}

type CreateFromTemplate struct {
	TemplateID  uint   `json:"template_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	GithubLink  string `json:"github_link"`
}

// ========== EXPORT ==========

type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportPDF  ExportFormat = "pdf"
)

type ExportRequest struct {
	Format       ExportFormat `json:"format" binding:"required,oneof=json pdf"`
	IncludeStats bool         `json:"include_stats"`
}

type ProjectExportData struct {
	Project      ProjectResponse   `json:"project"`
	Statistics   *ProjectStats     `json:"statistics,omitempty"`
	Comments     []CommentResponse `json:"comments,omitempty"`
	Reviews      []ReviewResponse  `json:"reviews,omitempty"`
	ExportedAt   time.Time         `json:"exported_at"`
	ExportFormat string            `json:"export_format"`
}

type PortfolioExport struct {
	ExportDate    time.Time          `json:"export_date"`
	TotalProjects int                `json:"total_projects"`
	Projects      []ProjectResponse  `json:"projects"`
	Statistics    *PlatformAnalytics `json:"statistics,omitempty"`
}

// ========== COMMENT & REVIEW ==========

type CommentCreate struct {
	ProjectID uint   `json:"project_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID        uint        `json:"id"`
	ProjectID uint        `json:"project_id"`
	Content   string      `json:"content"`
	User      Contributor `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}

type ReviewCreate struct {
	ProjectID uint    `json:"project_id" binding:"required"`
	Rating    float64 `json:"rating" binding:"required,min=1,max=5"`
	Comment   string  `json:"comment"`
}

type ReviewResponse struct {
	ID        uint        `json:"id"`
	ProjectID uint        `json:"project_id"`
	Rating    float64     `json:"rating"`
	Comment   string      `json:"comment"`
	User      Contributor `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}

type ProjectStats struct {
	ViewCount     int     `json:"view_count"`
	LikeCount     int     `json:"like_count"`
	CommentCount  int     `json:"comment_count"`
	ReviewCount   int     `json:"review_count"`
	AverageRating float64 `json:"average_rating"`
}

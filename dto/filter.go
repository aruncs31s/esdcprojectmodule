package dto

type ProjectFilter struct {
	Category     string   `form:"category"`
	Technologies []string `form:"technologies"`
	Status       string   `form:"status"`
	Visibility   string   `form:"visibility"`
	Search       string   `form:"search"`
	Limit        int      `form:"limit" binding:"min=1,max=100"`
	Offset       int      `form:"offset" binding:"min=0"`
}

type PaginationMetadata struct {
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
}

type ProjectListResponse struct {
	Projects   []ProjectResponse   `json:"projects"`
	Pagination PaginationMetadata  `json:"pagination"`
}

type ProjectUpdate struct {
	Title        *string   `json:"title"`
	Description  *string   `json:"description"`
	Image        *string   `json:"image"`
	Status       *string   `json:"status"`
	Visibility   *string   `json:"visibility"`
	GithubLink   *string   `json:"github_link"`
	LiveURL      *string   `json:"live_url"`
	Category     *string   `json:"category"`
	Technologies *[]string `json:"technologies"`
	Tags         *[]string `json:"tags"`
}

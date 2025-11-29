package dto

import "time"

type CommentCreate struct {
	Content   string `json:"content" binding:"required,min=1,max=1000"`
	ProjectID int    `json:"project_id" binding:"required"`
}

type CommentResponse struct {
	ID        uint        `json:"id"`
	Content   string      `json:"content"`
	UserID    uint        `json:"user_id"`
	ProjectID uint        `json:"project_id"`
	User      Contributor `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type ReviewCreate struct {
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"max=500"`
	ProjectID int    `json:"project_id" binding:"required"`
}

type ReviewResponse struct {
	ID        uint        `json:"id"`
	Rating    int         `json:"rating"`
	Comment   string      `json:"comment"`
	UserID    uint        `json:"user_id"`
	ProjectID uint        `json:"project_id"`
	User      Contributor `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
}

type ProjectStats struct {
	ViewCount     int     `json:"view_count"`
	LikeCount     int     `json:"like_count"`
	CommentCount  int     `json:"comment_count"`
	ForkCount     int     `json:"fork_count"`
	AverageRating float64 `json:"average_rating"`
	ReviewCount   int     `json:"review_count"`
}

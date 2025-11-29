package repository

import (
	"time"

	commonModules "github.com/aruncs31s/esdcmodels"
)

// GetTrendingProjects returns trending projects based on engagement metrics
// Trending score = (likes * 0.4) + (views * 0.3) + (comments * 0.2) + (reviews * 0.1)
func (r *projectRepositoryReader) GetTrendingProjects(limit, offset int, days int) ([]commonModules.Project, error) {
	var projects []commonModules.Project

	// Calculate trending score for the last N days
	since := time.Now().AddDate(0, 0, -days)

	err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("created_at >= ?", since).
		Order("(likes * 0.4 + views * 0.3) DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, err
}

// GetRecommendedProjects returns personalized recommendations based on user's project interests
// Recommendations based on:
// 1. Projects with similar technologies
// 2. Projects in the same category
// 3. Popular projects by users they follow
func (r *projectRepositoryReader) GetRecommendedProjects(userID uint, limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project

	// Get user's projects to understand their interests
	var userTechs []uint
	r.db.Raw(`
		SELECT DISTINCT pt.technology_id 
		FROM project_technologies pt
		JOIN projects p ON p.id = pt.project_id
		WHERE p.created_by = ?
	`, userID).Scan(&userTechs)

	// Find projects with similar technologies (excluding user's own projects)
	query := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("created_by != ?", userID).
		Where("visibility = ?", 0)

	if len(userTechs) > 0 {
		query = query.
			Joins("JOIN project_technologies pt ON projects.id = pt.project_id").
			Where("pt.technology_id IN ?", userTechs)
	}

	err := query.
		Order("likes DESC, views DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, err
}

// GetSimilarProjects returns projects similar to the given project
// Similarity based on: shared technologies, same category, overlapping tags
func (r *projectRepositoryReader) GetSimilarProjects(projectID uint, limit, offset int) ([]commonModules.Project, error) {
	var projects []commonModules.Project

	// Get the reference project's technologies and category
	var refProject commonModules.Project
	if err := r.db.First(&refProject, projectID).Error; err != nil {
		return nil, err
	}

	// Find projects with matching technologies or category
	err := r.db.
		Preload("Contributors").
		Preload("Creator").
		Preload("Tags").
		Preload("Technologies").
		Where("id != ? AND visibility = ?", projectID, 0).
		Where("category = ? OR id IN (SELECT project_id FROM project_technologies WHERE technology_id IN (SELECT technology_id FROM project_technologies WHERE project_id = ?))", refProject.Category, projectID).
		Order("likes DESC, views DESC").
		Limit(limit).
		Offset(offset).
		Find(&projects).Error

	return projects, err
}

// GetProjectAnalytics returns analytics for a specific project
func (r *projectRepositoryReader) GetProjectAnalytics(projectID uint, days int) (*commonModules.ProjectStats, error) {
	var stats commonModules.ProjectStats
	var project commonModules.Project

	if err := r.db.First(&project, projectID).Error; err != nil {
		return nil, err
	}

	since := time.Now().AddDate(0, 0, -days)

	var commentCount, reviewCount int64
	var avgRating float64

	r.db.Model(&commonModules.Comment{}).
		Where("project_id = ? AND created_at >= ?", projectID, since).
		Count(&commentCount)

	r.db.Model(&commonModules.Review{}).
		Where("project_id = ? AND created_at >= ?", projectID, since).
		Count(&reviewCount)

	r.db.Model(&commonModules.Review{}).
		Where("project_id = ? AND created_at >= ?", projectID, since).
		Select("COALESCE(AVG(rating), 0)").
		Scan(&avgRating)

	stats.ViewCount = project.Views
	stats.LikeCount = project.Likes
	stats.CommentCount = int(commentCount)
	stats.ReviewCount = int(reviewCount)
	stats.AverageRating = avgRating

	return &stats, nil
}

// GetPlatformAnalytics returns platform-wide analytics
func (r *projectRepositoryReader) GetPlatformAnalytics(days int) (*commonModules.PlatformAnalytics, error) {
	var analytics commonModules.PlatformAnalytics

	since := time.Now().AddDate(0, 0, -days)

	// Total projects
	r.db.Model(&commonModules.Project{}).Count((*int64)(&analytics.TotalProjects))

	// Total views and likes
	r.db.Model(&commonModules.Project{}).
		Select("COALESCE(SUM(views), 0) as total_views, COALESCE(SUM(likes), 0) as total_likes").
		Where("created_at >= ?", since).
		Row().
		Scan(&analytics.TotalViews, &analytics.TotalLikes)

	return &analytics, nil
}

// GetTrendingTechnologies returns the most popular technologies
func (r *projectRepositoryReader) GetTrendingTechnologies(limit int) ([]commonModules.TrendingTech, error) {
	var techs []commonModules.TrendingTech

	err := r.db.Raw(`
		SELECT t.id, t.name, COUNT(pt.project_id) as usage_count
		FROM technologies t
		LEFT JOIN project_technologies pt ON t.id = pt.technology_id
		GROUP BY t.id, t.name
		ORDER BY usage_count DESC
		LIMIT ?
	`, limit).Scan(&techs).Error

	return techs, err
}

// GetTrendingTags returns the most popular tags
func (r *projectRepositoryReader) GetTrendingTags(limit int) ([]commonModules.TrendingTag, error) {
	var tags []commonModules.TrendingTag

	err := r.db.Raw(`
		SELECT t.id, t.name, COUNT(pt.project_id) as usage_count
		FROM tags t
		LEFT JOIN project_tags pt ON t.id = pt.tag_id
		GROUP BY t.id, t.name
		ORDER BY usage_count DESC
		LIMIT ?
	`, limit).Scan(&tags).Error

	return tags, err
}

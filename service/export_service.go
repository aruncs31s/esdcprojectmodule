package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aruncs31s/esdcprojectmodule/dto"
	"github.com/aruncs31s/esdcprojectmodule/interfaces/repository"
	userRepo "github.com/aruncs31s/esdcusermodule/repository"
)

type ExportService interface {
	ExportProject(projectID uint, username string, format dto.ExportFormat) ([]byte, error)
	ExportPortfolio(username string, format dto.ExportFormat, includeStats bool) ([]byte, error)
}

type exportService struct {
	projectRepo repository.ProjectRepository
	userRepo    userRepo.UserRepository
}

func NewExportService(
	projectRepo repository.ProjectRepository,
	userRepo userRepo.UserRepository,
) ExportService {
	return &exportService{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// ExportProject exports a single project in the specified format
func (s *exportService) ExportProject(projectID uint, username string, format dto.ExportFormat) ([]byte, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	// Fetch project details
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Only project creator can export
	if project.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized to export this project")
	}

	// Get project stats
	stats, _ := s.projectRepo.GetProjectStats(projectID)

	// Get comments and reviews
	comments, _ := s.projectRepo.GetComments(projectID, 100, 0)
	reviews, _ := s.projectRepo.GetReviews(projectID, 100, 0)

	// Build export data
	exportData := dto.ProjectExportData{
		Project:      *getProjectResponseForPersonal(project, false),
		Statistics:   stats,
		ExportedAt:   time.Now(),
		ExportFormat: string(format),
	}

	// Convert comments
	for _, comment := range comments {
		exportData.Comments = append(exportData.Comments, dto.CommentResponse{
			ID:        comment.ID,
			ProjectID: comment.ProjectID,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		})
	}

	// Convert reviews
	for _, review := range reviews {
		exportData.Reviews = append(exportData.Reviews, dto.ReviewResponse{
			ID:        review.ID,
			ProjectID: review.ProjectID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
		})
	}

	switch format {
	case dto.ExportJSON:
		return s.exportToJSON(exportData)
	case dto.ExportPDF:
		return s.exportToPDF(exportData)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// ExportPortfolio exports all user projects
func (s *exportService) ExportPortfolio(username string, format dto.ExportFormat, includeStats bool) ([]byte, error) {
	userID, err := s.userRepo.FindUserIDByUsername(username)
	if err != nil {
		return nil, err
	}

	// Fetch all user projects
	projects, err := s.projectRepo.GetUserProjects(userID, 1000, 0)
	if err != nil {
		return nil, err
	}

	// Build portfolio export data
	var projectResponses []dto.ProjectResponse
	for _, project := range projects {
		projectResponses = append(projectResponses, *getProjectResponseForPersonal(project, false))
	}

	portfolio := dto.PortfolioExport{
		ExportDate:    time.Now(),
		TotalProjects: len(projectResponses),
		Projects:      projectResponses,
	}

	if includeStats {
		analytics, _ := s.projectRepo.GetPlatformAnalytics(30)
		portfolio.Statistics = (*dto.PlatformAnalytics)(analytics)
	}

	switch format {
	case dto.ExportJSON:
		return s.exportPortfolioToJSON(portfolio)
	case dto.ExportPDF:
		return s.exportPortfolioPDF(portfolio)
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}

// exportToJSON converts project export data to JSON
func (s *exportService) exportToJSON(data dto.ProjectExportData) ([]byte, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return jsonData, nil
}

// exportPortfolioToJSON converts portfolio export data to JSON
func (s *exportService) exportPortfolioToJSON(data dto.PortfolioExport) ([]byte, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return jsonData, nil
}

// exportToPDF converts project export data to PDF
// This uses a placeholder implementation - integrate with actual PDF library (e.g., go-pdf)
func (s *exportService) exportToPDF(data dto.ProjectExportData) ([]byte, error) {
	// For full implementation, use github.com/go-pdf/fpdf or similar
	// This is a placeholder that returns JSON wrapped as PDF content
	jsonData, _ := json.MarshalIndent(data, "", "  ")

	// In a real implementation, you would:
	// 1. Create a PDF document
	// 2. Add project information
	// 3. Add statistics
	// 4. Add comments and reviews
	// 5. Generate and return PDF bytes

	return jsonData, nil // Placeholder
}

// exportPortfolioPDF converts portfolio data to PDF format
func (s *exportService) exportPortfolioPDF(data dto.PortfolioExport) ([]byte, error) {
	// Placeholder implementation
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	return jsonData, nil
}

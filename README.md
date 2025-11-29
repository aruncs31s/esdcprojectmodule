# Project Module

A Go module for managing projects in the [ESDC Website](https://esdc.vercel.app/).

## Features

- **Project Management** - Create, update, delete, and list projects
- **Trending & Recommendations** - Trending projects, personalized recommendations, similar projects
- **Analytics** - Project analytics and platform-wide statistics
- **Templates** - Save projects as templates, template marketplace, create from templates
- **Notifications** - Like/comment notifications, follower updates, milestone notifications
- **Export** - Export projects and portfolios to JSON/PDF
- **Comments & Reviews** - User comments and ratings on projects

## Installation

```bash
go get github.com/aruncs31s/esdcprojectmodule
```

## Quick Start

```go
package main

import (
	project "github.com/aruncs31s/esdcprojectmodule"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	
	project.InitProjectModule(r, db)
	project.RegisterPublicProjectRoutes()
	project.RegisterProjectRoutes()  // Protected routes (requires auth)
	
	r.Run()
}
```

## API Endpoints

### Public Endpoints (No Authentication)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/public/projects` | List all public projects |
| GET | `/api/public/projects/:id` | Get project by ID |
| GET | `/api/public/projects/trending` | Get trending projects |
| GET | `/api/public/projects/:id/similar` | Get similar projects |
| GET | `/api/public/projects/templates` | List public templates |
| GET | `/api/public/projects/analytics/platform` | Platform-wide analytics |

### Protected Endpoints (Bearer Token Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects` | Create a new project |
| PUT | `/api/projects/:id` | Update a project |
| DELETE | `/api/projects/:id` | Delete a project |
| POST | `/api/projects/:id/like` | Like a project |
| DELETE | `/api/projects/:id/like` | Unlike a project |
| GET | `/api/projects/recommendations` | Get personalized recommendations |
| GET | `/api/projects/:id/analytics` | Get project analytics (owner only) |
| GET | `/api/projects/:id/export` | Export project (JSON/PDF) |
| GET | `/api/projects/portfolio/export` | Export user portfolio |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/projects/templates` | Create a template from project |
| GET | `/api/projects/templates/user` | Get user's templates |
| GET | `/api/projects/templates/:id` | Get template details |
| DELETE | `/api/projects/templates/:id` | Delete a template |
| POST | `/api/projects/from-template` | Create project from template |

### Notifications

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/projects/notifications` | Get user notifications |
| POST | `/api/projects/notifications/:id/read` | Mark notification as read |
| POST | `/api/projects/notifications/read-all` | Mark all as read |
| DELETE | `/api/projects/notifications/:id` | Delete notification |

### Comments & Reviews

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/projects/:id/comments` | Get project comments |
| POST | `/api/projects/:id/comments` | Add a comment |
| GET | `/api/projects/:id/reviews` | Get project reviews |
| POST | `/api/projects/:id/reviews` | Add a review |

## Dependencies

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [esdcmodels](https://github.com/aruncs31s/esdcmodels) - Shared data models
- [esdcusermodule](https://github.com/aruncs31s/esdcusermodule) - User management

## Project Structure

```
esdcprojectmodule/
├── dto/                    # Data Transfer Objects
│   ├── project.go          # Project DTOs (requests/responses)
│   └── comment.go          # Comment/Review DTOs
├── handler/                # HTTP Handlers
│   ├── projects_handler.go
│   └── public_project_handler.go
├── interfaces/             # Interface definitions
│   ├── handler/
│   ├── repository/
│   └── service/
├── repository/             # Database operations
│   ├── projects_repository.go
│   ├── trending_analytics_repository.go
│   ├── template_repository.go
│   └── notification_repository.go
├── routes/                 # Route registration
│   └── project_routes.go
├── service/                # Business logic
│   ├── project_service.go
│   ├── project_service_extended.go
│   ├── trending_recommendations_service.go
│   └── export_service.go
├── utils/                  # Utility functions
│   └── project_utils.go
├── project.go              # Module initialization
└── README.md
```

## Documentation

- [Advanced Features](./ADVANCED_FEATURES.md) - Detailed API documentation
- [API Quick Reference](./API_QUICK_REFERENCE.md) - Quick lookup with examples
- [Models Guide](./MODELS_TO_ADD.md) - Database model definitions

## License

See [LICENSE](./LICENSE) for details.
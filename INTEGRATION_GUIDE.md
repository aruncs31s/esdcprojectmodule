# ESDC Project Module - Integration & Testing Guide

## Quick Start Guide

### Initialization
```go
package main

import (
    project "github.com/aruncs31s/esdcprojectmodule"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func main() {
    r := gin.Default()
    db, _ := gorm.Open(sqlite.Open("db.db"), &gorm.Config{})
    
    // Initialize the module
    project.InitProjectModule(r, db)
    
    // Register all route groups
    project.RegisterPublicProjectRoutes()
    
    // Private routes need JWT middleware
    privateGroup := r.Group("/api/projects")
    privateGroup.Use(yourJWTMiddleware())
    project.RegisterPrivateProjectRoutes(r)
    
    // Admin routes need admin middleware
    adminGroup := r.Group("/api/admin/projects")
    adminGroup.Use(yourAdminMiddleware())
    project.RegisterAdminProjectRoutes(r)
    
    r.Run()
}
```

---

## API Testing Examples

### 1. Filtering & Search

#### Filter by Category
```bash
curl -X GET "http://localhost:8080/api/public/projects/filter?category=Web%20Development&limit=10&offset=0"
```

#### Filter by Technology
```bash
curl -X GET "http://localhost:8080/api/public/projects/filter?technologies=Go,React&limit=10&offset=0"
```

#### Full-Text Search
```bash
curl -X GET "http://localhost:8080/api/public/projects/filter?search=api&limit=10&offset=0"
```

#### Advanced Filtering (Multiple Criteria)
```bash
curl -X GET "http://localhost:8080/api/public/projects/filter?category=Web%20Development&status=active&visibility=public&search=api&limit=10&offset=0"
```

### Response
```json
{
  "projects": [
    {
      "id": 1,
      "title": "My API Project",
      "description": "RESTful API built with Go",
      "category": "Web Development",
      "status": "active",
      ...
    }
  ],
  "pagination": {
    "total": 42,
    "total_pages": 5,
    "current_page": 1,
    "per_page": 10
  }
}
```

---

### 2. Update Project

```bash
curl -X PUT "http://localhost:8080/api/projects/1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Project Title",
    "description": "Updated description",
    "status": "completed",
    "category": "Backend Development"
  }'
```

### Response
```json
{
  "message": "project updated successfully"
}
```

---

### 3. Delete Project

#### Soft Delete (Owner)
```bash
curl -X DELETE "http://localhost:8080/api/projects/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Hard Delete (Admin)
```bash
curl -X DELETE "http://localhost:8080/api/admin/projects/1" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Response
```json
{
  "message": "project deleted successfully"
}
```

---

### 4. Get Project Statistics

```bash
curl -X GET "http://localhost:8080/api/public/projects/1/stats"
```

### Response
```json
{
  "view_count": 1250,
  "like_count": 45,
  "comment_count": 12,
  "fork_count": 8,
  "average_rating": 4.5,
  "review_count": 10
}
```

---

### 5. Comments Management

#### Create Comment
```bash
curl -X POST "http://localhost:8080/api/projects/comments" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Great project! Very helpful.",
    "project_id": 1
  }'
```

#### Get Comments (Paginated)
```bash
curl -X GET "http://localhost:8080/api/public/projects/1/comments?limit=10&offset=0"
```

#### Delete Comment
```bash
curl -X DELETE "http://localhost:8080/api/projects/comments/1" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### Moderate Comment (Admin)
```bash
curl -X PATCH "http://localhost:8080/api/admin/projects/comments/1/moderate" \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved"
  }'
```

### Comment Response
```json
{
  "id": 1,
  "content": "Great project!",
  "user_id": 5,
  "project_id": 1,
  "user": {
    "id": 5,
    "name": "John Doe",
    "email": "john@example.com",
    "image": "https://avatar.jpg"
  },
  "created_at": "2025-01-15T10:30:00Z",
  "updated_at": "2025-01-15T10:30:00Z"
}
```

---

### 6. Reviews Management

#### Create Review
```bash
curl -X POST "http://localhost:8080/api/projects/reviews" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 5,
    "comment": "Excellent implementation!",
    "project_id": 1
  }'
```

#### Get Reviews (Paginated)
```bash
curl -X GET "http://localhost:8080/api/public/projects/1/reviews?limit=10&offset=0"
```

### Review Response
```json
{
  "id": 1,
  "rating": 5,
  "comment": "Excellent implementation!",
  "user_id": 5,
  "project_id": 1,
  "user": {
    "id": 5,
    "name": "John Doe",
    "email": "john@example.com",
    "image": "https://avatar.jpg"
  },
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### 7. Like/Unlike Project

```bash
curl -X POST "http://localhost:8080/api/projects/1/toggle-like" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response
```json
{
  "liked": true,
  "message": "Like toggled successfully"
}
```

---

## Middleware Requirements

### JWT Middleware (For Private Routes)
The middleware should:
1. Extract JWT token from Authorization header
2. Validate the token
3. Extract username and set in context as `username`
4. Set `isAdmin` if applicable

Example:
```go
func JWTMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // Validate token
        username := extractUsername(token)
        c.Set("username", username)
        c.Set("isAdmin", isUserAdmin(username))
        c.Next()
    }
}
```

### Admin Middleware (For Admin Routes)
```go
func AdminMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        isAdmin, _ := c.Get("isAdmin")
        if !isAdmin.(bool) {
            c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

---

## Error Handling

### Common Error Responses

#### 400 Bad Request
```json
{
  "error": "invalid project ID"
}
```

#### 401 Unauthorized
```json
{
  "error": "unauthorized: only project owner can update"
}
```

#### 404 Not Found
```json
{
  "error": "Project not found"
}
```

#### 500 Internal Server Error
```json
{
  "error": "Failed to create project"
}
```

---

## Testing Checklist

### Unit Tests
- [ ] Test FilterProjects with various criteria combinations
- [ ] Test UpdateProject authorization
- [ ] Test DeleteProject (soft vs hard)
- [ ] Test GetProjectStats calculations
- [ ] Test CreateComment validation
- [ ] Test CreateReview rating validation
- [ ] Test ToggleLikeProject idempotence
- [ ] Test pagination calculations

### Integration Tests
- [ ] Test full filter + pagination flow
- [ ] Test comment moderation workflow
- [ ] Test review creation and rating calculation
- [ ] Test like/unlike operations
- [ ] Test project update with all fields
- [ ] Test soft delete preservation
- [ ] Test hard delete removal

### Authorization Tests
- [ ] Non-owner cannot update project
- [ ] Non-owner cannot soft delete project
- [ ] Non-admin cannot hard delete project
- [ ] Non-admin cannot moderate comments
- [ ] Public routes accessible without auth
- [ ] Private routes require auth

### Edge Cases
- [ ] Empty search results
- [ ] Pagination beyond available data
- [ ] Invalid filter criteria
- [ ] Update with empty/null fields
- [ ] Delete already deleted projects
- [ ] Rating outside 1-5 range
- [ ] Comment exceeding 1000 characters

---

## Performance Optimization Tips

1. **Database Indexing**
   ```sql
   CREATE INDEX idx_project_category ON projects(category);
   CREATE INDEX idx_project_status ON projects(status);
   CREATE INDEX idx_project_visibility ON projects(visibility);
   CREATE INDEX idx_comment_project ON comments(project_id);
   CREATE INDEX idx_review_project ON reviews(project_id);
   ```

2. **Query Optimization**
   - Use Preload for relations (already implemented)
   - Limit results with pagination (already implemented)
   - Use specific SELECT for large queries

3. **Caching Recommendations**
   - Cache statistics for popular projects
   - Cache filtered results for common queries
   - Cache category/technology lists

4. **Rate Limiting**
   - Limit filter requests per user/IP
   - Limit comment creation
   - Limit review creation

---

## Database Schema Requirements

Ensure the following tables exist in your database:

### projects
```sql
CREATE TABLE projects (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  category VARCHAR(100),
  status VARCHAR(50),
  visibility INT,
  github_link VARCHAR(255),
  live_url VARCHAR(255),
  image VARCHAR(255),
  views INT DEFAULT 0,
  likes INT DEFAULT 0,
  cost INT,
  created_by BIGINT,
  modified_by BIGINT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  FOREIGN KEY (created_by) REFERENCES users(id),
  FOREIGN KEY (modified_by) REFERENCES users(id)
);
```

### comments
```sql
CREATE TABLE comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  content TEXT NOT NULL,
  user_id BIGINT NOT NULL,
  project_id BIGINT NOT NULL,
  status VARCHAR(50) DEFAULT 'approved',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

### reviews
```sql
CREATE TABLE reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
  comment VARCHAR(500),
  user_id BIGINT NOT NULL,
  project_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

### project_likes (for like tracking)
```sql
CREATE TABLE project_likes (
  user_id BIGINT NOT NULL,
  project_id BIGINT NOT NULL,
  PRIMARY KEY (user_id, project_id),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (project_id) REFERENCES projects(id)
);
```

### technologies
```sql
CREATE TABLE technologies (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) UNIQUE NOT NULL
);
```

### project_technologies
```sql
CREATE TABLE project_technologies (
  project_id BIGINT NOT NULL,
  technology_id BIGINT NOT NULL,
  PRIMARY KEY (project_id, technology_id),
  FOREIGN KEY (project_id) REFERENCES projects(id),
  FOREIGN KEY (technology_id) REFERENCES technologies(id)
);
```

---

## Monitoring & Logging

### Recommended Logging Points
1. Filter operations (searches)
2. Update/Delete operations
3. Authorization failures
4. Database errors
5. Validation errors
6. Performance metrics

### Metrics to Track
- Requests per endpoint
- Average response time
- Error rates
- Most filtered categories
- Most searched terms
- Popular projects (by view count)

---

## Future Enhancements

1. **Advanced Search**
   - Fuzzy search for typos
   - Search suggestions
   - Popular searches

2. **Recommendations**
   - Similar projects
   - Recommended projects based on views
   - Trending projects

3. **Analytics**
   - Usage analytics dashboard
   - Search analytics
   - Popular filters

4. **Social Features**
   - Follow projects
   - Share projects
   - Notifications

5. **Batch Operations**
   - Bulk update status
   - Bulk delete
   - Bulk export

---

## Support & Troubleshooting

### Common Issues

**Issue: Filter returns no results**
- Check if projects exist with matching criteria
- Verify visibility settings
- Check technology exact spelling

**Issue: Update fails with 401**
- Verify user is project owner
- Check JWT token validity
- Ensure username is set in context

**Issue: Comments not appearing**
- Check comment status (approved/rejected/pending)
- Verify project exists
- Check pagination parameters

**Issue: Statistics showing 0**
- Ensure IncrementViewCount is called
- Check database has related records
- Verify aggregation queries

---

## Dependencies

The module requires:
- `github.com/aruncs31s/esdcmodels` - Data models
- `github.com/aruncs31s/esdcusermodule` - User management
- `github.com/aruncs31s/esdcsharedhelpersmodule` - Helper utilities
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` - ORM

---

## Version Information

- **Module Version**: v1.0.0
- **Go Version**: 1.16+
- **Gin Version**: 1.7+
- **GORM Version**: 1.23+


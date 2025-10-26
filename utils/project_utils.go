package utils

import (
	commonModules "github.com/aruncs31s/esdcmodels"
	"github.com/aruncs31s/esdcprojectmodule/dto"
)

// Move them to a common group.
func GetCreatorDetails(creator commonModules.User) dto.Contributor {
	return dto.Contributor{
		ID:    int(creator.ID),
		Name:  creator.Username,
		Email: creator.Email,
		Image: getImage(creator.Image),
	}
}
func getImage(image *string) string {
	if image != nil {
		return *image
	}
	// Default image URL
	return "https://example.com/default-image.png"
}
func GetContributorsUsernames(contributors *[]commonModules.User) *[]dto.Contributor {
	if contributors == nil {
		return nil
	}
	usernames := make([]dto.Contributor, 0)
	for _, user := range *contributors {
		usernames = append(usernames, dto.Contributor{
			ID:    int(user.ID),
			Name:  user.Username,
			Image: getImage(user.Image),
			Email: user.Email,
		})
	}
	// add creator as a contributor also if not already present
	return &usernames
}
func GetTagsNames(tags *[]commonModules.Tag) *[]dto.Tag {
	if tags == nil {
		return nil
	}
	names := make([]dto.Tag, 0)
	for _, tag := range *tags {
		names = append(names,
			dto.Tag{
				ID:   int(tag.ID),
				Name: tag.Name,
			})
	}
	return &names
}
func GetTechnologiesNames(technologies *[]commonModules.Technologies) *[]dto.Technology {
	if technologies == nil {
		return nil
	}
	names := make([]dto.Technology, 0)

	for _, tech := range *technologies {
		names = append(names, dto.Technology{
			ID:   int(tech.ID),
			Name: tech.Name,
		})
	}
	return &names
}

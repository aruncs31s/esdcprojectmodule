package repository

import commonModules "github.com/aruncs31s/esdcmodels"

func (r *projectRepositoryReader) GetEssentialInfo(limit, offset int) (*[]commonModules.Project, error) {
	var projects []commonModules.Project
	essentialFields := commonModules.Project{}.GetProjectEssentialFields()
	query := r.db.
		Model(&commonModules.Project{}).
		Select(essentialFields).
		Limit(limit).Offset(offset)
	if err := query.Find(&projects).Error; err != nil {
		return nil, err
	}
	return &projects, nil
}

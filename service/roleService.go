package service

import (
	"srbolab_cpc/db"
	"srbolab_cpc/model"
)

var (
	RoleService roleServiceInterface = &roleService{}
)

type roleService struct {
}

type roleServiceInterface interface {
	GetAllRoles() ([]model.Role, error)
	GetRoleByID(id int) (*model.Role, error)
}

func (s *roleService) GetAllRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := db.Client.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *roleService) GetRoleByID(id int) (*model.Role, error) {
	var role *model.Role
	if err := db.Client.First(&role, id).Error; err != nil {
		return nil, err
	}

	return role, nil
}

package repository

import (
	"project-go/config"
	"project-go/internal/model"

	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id string) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id string) error
	GetAll() ([]model.Role, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() RoleRepository {
	return &roleRepository{
		db: config.DB,
	}
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id string) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, "id = ?", id).Error
	return &role, err
}

func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id string) error {
	return r.db.Delete(&model.Role{}, "id = ?", id).Error
}

func (r *roleRepository) GetAll() ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Find(&roles).Error
	return roles, err
}
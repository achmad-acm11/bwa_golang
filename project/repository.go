package project

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(project Project) (Project, error)
	Update(project Project) (Project, error)
	GetAll() ([]Project, error)
	GetAllByUser(id_user int) ([]Project, error)
	GetOneById(id int) (Project, error)
	GetOneByUser(id_user int) (Project, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Query Create Project
func (r *repository) Create(project Project) (Project, error) {
	err := r.db.Create(&project).Error

	if err != nil {
		return project, err
	}

	return project, nil
}

// Query Update Project
func (r *repository) Update(project Project) (Project, error) {
	err := r.db.Save(&project).Error

	if err != nil {
		return project, err
	}

	return project, nil
}

// Query List Project All
func (r *repository) GetAll() ([]Project, error) {
	project := []Project{}
	err := r.db.Preload("User").Preload("Images", "images.is_primary = ?", 1).Find(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

// Get List Project By User_id
func (r *repository) GetAllByUser(id_user int) ([]Project, error) {
	project := []Project{}
	err := r.db.Where("id_user = ?", id_user).Preload("Images", "images.is_primary = ?", 1).Find(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

// Query Get Data Project by Id
func (r *repository) GetOneById(id int) (Project, error) {
	project := Project{}
	err := r.db.Where("id = ?", id).Preload("User").Preload("Images").Find(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

// Query Get Data Project by Id_user
func (r *repository) GetOneByUser(id_user int) (Project, error) {
	project := Project{}
	err := r.db.Where("id_user = ?", id_user).Preload("images.is_primary = ?", 1).Find(&project).Error
	if err != nil {
		return project, err
	}
	return project, nil
}

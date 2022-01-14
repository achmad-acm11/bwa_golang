package image

import "gorm.io/gorm"

type Repository interface {
	CreateImage(image Image) (Image, error)
	ClearPrimaryImage(id_project int) error
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Query Create Image
func (r *repository) CreateImage(image Image) (Image, error) {
	err := r.db.Create(&image).Error

	if err != nil {
		return image, err
	}

	return image, nil
}

// Query Clear Primary Image
func (r *repository) ClearPrimaryImage(id_project int) error {
	err := r.db.Model(&Image{}).Where("id_project = ?", id_project).Update("is_primary", 0).Error
	if err != nil {
		return err
	}
	return nil
}

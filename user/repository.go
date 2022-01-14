package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(user User) (User, error)
	FindById(id int) (User, error)
	Update(user User) (User, error)
	GetAll() ([]User, error)
	Delete(id int) error
}
type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Query Get All Data
func (r *repository) GetAll() ([]User, error) {
	users := []User{}
	err := r.db.Find(&users).Error

	if err != nil {
		return users, err
	}

	return users, nil
}

// Query Create Data
func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// Query Delete Data
func (r *repository) Delete(id int) error {
	err := r.db.Delete(User{}, "id = ?", id).Error

	if err != nil {
		return err
	}
	return nil
}

// Query Update User
func (r *repository) Update(user User) (User, error) {
	err := r.db.Save(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// Query Get User By Email
func (r *repository) FindByEmail(user User) (User, error) {
	var userRes User

	err := r.db.Where("email = ?", user.Email).Find(&userRes).Error

	if err != nil {
		return userRes, err
	}

	return userRes, nil
}

// Query Get Data By Id
func (r *repository) FindById(id int) (User, error) {
	var userRes User

	err := r.db.Where("id = ?", id).Find(&userRes).Error

	if err != nil {
		return userRes, err
	}

	return userRes, nil
}

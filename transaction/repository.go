package transaction

import "gorm.io/gorm"

type Repository interface {
	GetAll() ([]Transaction, error)
	GetTransactionByProjectId(id_project int) ([]Transaction, error)
	GetTransactionById(id int) (Transaction, error)
	GetTransactionByUserId(id_user int) ([]Transaction, error)
	Create(transaction Transaction) (Transaction, error)
	Update(transaction Transaction) (Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Query Get List Transaction by Id_project
func (r *repository) GetAll() ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Preload("Project.Images", "images.is_primary = ?", 1).Order("id desc").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// Query Get List Transaction by Id_project
func (r *repository) GetTransactionByProjectId(id_project int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Where("id_project = ?", id_project).Order("id desc").Find(&transaction).Error

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// Query Get List Transaction by User
func (r *repository) GetTransactionByUserId(id_user int) ([]Transaction, error) {
	var transaction []Transaction
	err := r.db.Preload("User").Preload("Project.Images", "images.is_primary = ?", 1).Where("id_user = ?", id_user).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

// Query Get Data Transaction by Id_transaction
func (r *repository) GetTransactionById(id int) (Transaction, error) {
	var transaction Transaction
	err := r.db.Where("id = ?", id).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// Query Create Transaction
func (r *repository) Create(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// Query Update Transaction
func (r *repository) Update(transaction Transaction) (Transaction, error) {
	err := r.db.Save(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

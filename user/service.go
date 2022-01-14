package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	GetAllUser() ([]User, error)
	UpdateUser(user User) (User, error)
	LoginUser(input LoginUserInput) (User, error)
	CheckEmail(input EmailInput) (bool, error)
	GetUserById(id int) (User, error)
	UploadPhoto(id int, path string) (User, error)
	DeleteUser(id int) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

// Get User Data Service
func (s *service) GetAllUser() ([]User, error) {
	users, err := s.repository.GetAll()

	if err != nil {
		return users, err
	}

	return users, nil
}

// Create User Data Service
func (s *service) RegisterUser(i RegisterUserInput) (User, error) {
	// Generate Password
	passHash, err := bcrypt.GenerateFromPassword([]byte(i.Password), bcrypt.MinCost)

	if err != nil {
		return User{}, err
	}

	user := User{
		Name:       i.Name,
		Email:      i.Email,
		Password:   string(passHash),
		Profession: i.Profession,
		Role:       0,
	}
	// Create Data User Repo
	newUser, err := s.repository.Save(user)

	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

// Update User Data Service
func (s *service) UpdateUser(user User) (User, error) {
	dataUser, err := s.repository.Update(user)
	if err != nil {
		return dataUser, err
	}
	return dataUser, nil
}

// Login User Data Service
func (s *service) LoginUser(i LoginUserInput) (User, error) {
	user := User{
		Email: i.Email,
	}
	// Check Email User
	dataUser, err := s.repository.FindByEmail(user)

	if err != nil {
		return dataUser, err
	}
	// Data User is not exist
	if dataUser.Id == 0 {
		return dataUser, errors.New("no user found on that email")
	}
	// Check Password Hash
	err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(i.Password))
	if err != nil {
		return dataUser, err
	}

	return dataUser, nil
}

// Check Email User Service
func (s *service) CheckEmail(i EmailInput) (bool, error) {
	user := User{
		Email: i.Email,
	}
	// Get Data Email User
	dataUser, err := s.repository.FindByEmail(user)

	if err != nil {
		return false, err
	}
	// User Data is not Exist
	if dataUser.Id == 0 {
		return true, nil
	}

	return false, nil
}

// Get Data User By Id
func (s *service) GetUserById(id int) (User, error) {
	dataUser, err := s.repository.FindById(id)

	if err != nil {
		return dataUser, err
	}
	if dataUser.Id == 0 {
		return dataUser, errors.New("no user found on with that id")
	}
	return dataUser, nil
}

// Upload Photo User Data Service
func (s *service) UploadPhoto(id int, path string) (User, error) {
	// Get Data User Current
	dataUser, err := s.repository.FindById(id)
	if err != nil {
		return dataUser, err
	}
	// Update Data User
	dataUser.Image = path
	dataUser, err = s.repository.Update(dataUser)

	if err != nil {
		return dataUser, err
	}

	return dataUser, nil
}

func (s *service) DeleteUser(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

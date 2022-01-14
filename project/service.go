package project

import "errors"

type Service interface {
	GetAllProject() ([]Project, error)
	CreateProject(input CreateProjectInput) (Project, error)
	UpdateProject(input CreateProjectInput, currentProject Project) (Project, error)
	GetProjectById(id int) (Project, error)
	GetProjectByUser(id_user int) ([]Project, error)
}

type service struct {
	projectRepository *repository
}

func NewService(r *repository) *service {
	return &service{r}
}

// Get All Data Project
func (s *service) GetAllProject() ([]Project, error) {
	projects, err := s.projectRepository.GetAll()
	if err != nil {
		return projects, err
	}
	return projects, nil
}

// Create Project Service
func (s *service) CreateProject(input CreateProjectInput) (Project, error) {
	// Mapping Data Input to Project Model
	project := Project{
		Project_name:      input.Project_name,
		Description:       input.Description,
		Short_description: input.Short_description,
		Perks:             input.Perks,
		Slug:              input.Slug,
		Goal_amount:       input.Goal_amount,
		Current_amount:    0,
		Backer_count:      0,
		Id_user:           input.User.Id,
	}
	// Create Project
	project, err := s.projectRepository.Create(project)

	if err != nil {
		return project, err
	}

	return project, nil
}

// Update Project Service
func (s *service) UpdateProject(input CreateProjectInput, currentProject Project) (Project, error) {
	// Mapping Input Data to Project Model
	currentProject.Project_name = input.Project_name
	currentProject.Description = input.Description
	currentProject.Short_description = input.Short_description
	currentProject.Perks = input.Perks
	currentProject.Goal_amount = input.Goal_amount
	// Update Project
	project, err := s.projectRepository.Update(currentProject)

	if err != nil {
		return project, err
	}

	return project, nil
}

// Get Data Project By Id Service
func (s *service) GetProjectById(id int) (Project, error) {
	// Get Data Project
	project, err := s.projectRepository.GetOneById(id)

	if err != nil {
		return project, err
	}
	if project.Id == 0 {
		return project, errors.New("project not found")
	}
	return project, nil
}

// Get Data List Project by User_id or all Service
func (s *service) GetProjectByUser(id_user int) ([]Project, error) {
	// Get Data by user_id
	if id_user != 0 {
		project, err := s.projectRepository.GetAllByUser(id_user)
		if err != nil {
			return project, err
		}
		// User_id have no Project
		// if len(project) == 0 {
		// 	return project, errors.New("project not found")
		// }
		return project, nil
	}
	// Get Data List Project All
	project, err := s.projectRepository.GetAll()

	if err != nil {
		return project, err
	}
	// Empty Project
	if len(project) == 0 {
		return project, errors.New("project not found")
	}
	return project, nil
}

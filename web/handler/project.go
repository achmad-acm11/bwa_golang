package handler

import (
	"bwa_golang/image"
	"bwa_golang/project"
	"bwa_golang/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type projectHandler struct {
	projectService project.Service
	userService    user.Service
	imageService   image.Service
}

func NewProjectHandler(project project.Service, userService user.Service, imageService image.Service) *projectHandler {
	return &projectHandler{project, userService, imageService}
}

func (p *projectHandler) Index(c *gin.Context) {
	projects, err := p.projectService.GetAllProject()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "index_project.html", gin.H{"projects": projects})
}

func (p *projectHandler) New(c *gin.Context) {
	users, err := p.userService.GetAllUser()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input := project.CreateProjectAdmin{}
	input.Users = users

	c.HTML(http.StatusOK, "create_project.html", input)
}

func (p *projectHandler) Create(c *gin.Context) {
	var inputProject project.CreateProjectAdmin

	err := c.ShouldBind(&inputProject)
	if err != nil {
		users, e := p.userService.GetAllUser()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		inputProject.Users = users
		inputProject.Error = err.Error()
		c.HTML(http.StatusUnprocessableEntity, "create_project.html", inputProject)
		return
	}

	input := project.CreateProjectInput{
		Project_name:      inputProject.Project_name,
		Description:       inputProject.Description,
		Short_description: inputProject.Short_description,
		Perks:             inputProject.Perks,
		Goal_amount:       inputProject.Goal_amount,
	}
	// Generate Slug Project
	input.Slug = slug.Make(fmt.Sprintf("%s %d", input.Project_name, inputProject.UserId))
	dataUser, err := p.userService.GetUserById(inputProject.UserId)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	input.User = dataUser
	_, err = p.projectService.CreateProject(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/admin_project")
}

func (p *projectHandler) Image(c *gin.Context) {
	id := c.Param("id")
	id_project, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "create_image.html", gin.H{"id": id_project})
}
func (p *projectHandler) CreateImage(c *gin.Context) {
	id := c.Param("id")
	id_project, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	// Get Data File Image
	file, _ := c.FormFile("image")

	// Create Path Image Project
	path := fmt.Sprintf("project_image/%d-%s", id_project, file.Filename)
	// Save File Image to path
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	// Format Input Data Image & Validation
	var inputImage image.CreateImageInput
	inputImage.Image = path
	inputImage.Is_primary = true

	// Create Image service
	_, err = p.imageService.CreateImage(inputImage, id_project)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	c.Redirect(http.StatusFound, "/admin_project")
}

func (p *projectHandler) Edit(c *gin.Context) {
	// var uri_id user.GetUserIdUri
	id := c.Param("id")
	id_project, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	project, err := p.projectService.GetProjectById(id_project)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "edit_project.html", project)
}

func (p *projectHandler) Update(c *gin.Context) {
	var inputProject project.CreateProjectAdmin

	err := c.ShouldBind(&inputProject)

	if err != nil {
		inputProject.Error = err.Error()
		fmt.Println(inputProject.Error)
		// c.HTML(http.StatusUnprocessableEntity, "edit_user.html", nil)
		return
	}

	id := c.Param("id")
	id_project, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	dataProject, err := p.projectService.GetProjectById(id_project)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := project.CreateProjectInput{
		Project_name:      inputProject.Project_name,
		Description:       inputProject.Description,
		Short_description: inputProject.Short_description,
		Perks:             inputProject.Perks,
		Goal_amount:       inputProject.Goal_amount,
	}

	_, err = p.projectService.UpdateProject(input, dataProject)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/admin_project")
}

func (p *projectHandler) Show(c *gin.Context) {
	id := c.Param("id")
	id_project, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	project, err := p.projectService.GetProjectById(id_project)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "show_project.html", project)
}

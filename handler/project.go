package handler

import (
	"bwa_golang/helper"
	"bwa_golang/project"
	"bwa_golang/user"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type projectHandler struct {
	projectService project.Service
}

func NewProjectHandler(p project.Service) *projectHandler {
	return &projectHandler{p}
}

// Handler for Get List Project by user_id or all
func (p *projectHandler) GetProject(c *gin.Context) {
	// Get Data Id_user
	id_user, _ := strconv.Atoi(c.Query("id_user"))
	// Get Data List Project
	dataProject, err := p.projectService.GetProjectByUser(id_user)

	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success Get Project", "success", http.StatusOK, project.FormatterAll(dataProject))
	c.JSON(http.StatusOK, response)
}

// Handler for Get Detail Project
func (p *projectHandler) DetailProject(c *gin.Context) {
	// Get Id_project
	var uri_id project.GetProjectIdUri

	err := c.ShouldBindUri(&uri_id)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Data Project by Id_project service
	dataProject, err := p.projectService.GetProjectById(uri_id.Id)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Success Get Detail Project", "success", http.StatusOK, project.FormatterDetail(dataProject))
	c.JSON(http.StatusOK, response)
}

// Handler for Create Project
func (p *projectHandler) CreateProject(c *gin.Context) {
	// Format Input Data Project & Validation
	var input project.CreateProjectInput
	// Input Data For JSON
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Get Data Current User
	dataUser := c.MustGet("currentUser").(user.User)
	input.User = dataUser
	// Generate Slug Project
	input.Slug = slug.Make(fmt.Sprintf("%s %d", input.Project_name, dataUser.Id))
	// Create Project Service
	dataProject, err := p.projectService.CreateProject(input)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success Create Project", "success", http.StatusOK, project.Formatter(dataProject))
	c.JSON(http.StatusOK, response)
}

// Handler for Update Project
func (p *projectHandler) UpdateProject(c *gin.Context) {
	// Get Data id_project from URI
	var uri_id project.GetProjectIdUri
	err := c.ShouldBindUri(&uri_id)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Format Input Data & Validation
	var input project.CreateProjectInput
	// Get Data from JSON
	err = c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Get Data Project service
	dataProject, err := p.projectService.GetProjectById(uri_id.Id)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Data Project Not Found
	if dataProject.Id == 0 {
		response := helper.APIResponse("Project Not Found", "error", http.StatusBadRequest, gin.H{"errors": errors.New("project not found")})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Data User Id Login
	id := c.MustGet("currentUser").(user.User).Id

	if id != dataProject.Id_user {
		response := helper.APIResponse("You have not permission update", "error", http.StatusBadRequest, gin.H{"errors": errors.New("you have not permission update")})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Update Data Service
	dataProject, err = p.projectService.UpdateProject(input, dataProject)

	if err != nil {
		response := helper.APIResponse("Server Internal Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Success Update Project", "success", http.StatusOK, project.Formatter(dataProject))
	c.JSON(http.StatusOK, response)
}

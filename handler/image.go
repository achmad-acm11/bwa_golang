package handler

import (
	"bwa_golang/helper"
	"bwa_golang/image"
	"bwa_golang/project"
	"bwa_golang/user"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type imageHandler struct {
	imageSevice    image.Service
	projectService project.Service
}

func NewImageHandler(i image.Service, p project.Service) *imageHandler {
	return &imageHandler{i, p}
}

// Handlet for Upload Image Project
func (i *imageHandler) UploadImage(c *gin.Context) {
	// Get id_project
	var inputUri image.GetImageIdUri
	err := c.ShouldBindUri(&inputUri)

	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Project By Id_project service
	dataProject, err := i.projectService.GetProjectById(inputUri.Id)
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Id_user Login
	id_user := c.MustGet("currentUser").(user.User).Id

	if id_user != dataProject.Id_user {
		response := helper.APIResponse("You have not permission upload image", "error", http.StatusBadRequest, gin.H{"errors": errors.New("you have not permission upload")})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Data File Image
	file, _ := c.FormFile("image")

	// Create Path Image Project
	path := fmt.Sprintf("project_image/%d-%s", inputUri.Id, file.Filename)
	// Save File Image to path
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		response := helper.APIResponse("Upload File Failed", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Format Input Data Image & Validation
	var inputImage image.CreateImageInput
	// Get Input Data by form
	err = c.ShouldBind(&inputImage)

	inputImage.Image = path
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Create Image service
	dataImage, err := i.imageSevice.CreateImage(inputImage, inputUri.Id)
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadGateway, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Create Image Success", "success", http.StatusOK, image.FormatterImage(dataImage))
	c.JSON(http.StatusOK, response)
}

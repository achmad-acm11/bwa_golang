package handler

import (
	"bwa_golang/auth"
	"bwa_golang/helper"
	"bwa_golang/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(u user.Service, a auth.Service) *userHandler {
	return &userHandler{u, a}
}

// Handler for User Register
func (u *userHandler) RegisterUser(c *gin.Context) {
	// Format Input Data & Validasi
	var input user.RegisterUserInput
	// Get Input Data from type JSON
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Create User Data Service
	newUser, err := u.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Generate Token for auth JWT
	token, err := u.authService.GenerateToken(newUser.Id)

	if err != nil {
		response := helper.APIResponse("Failed generate token", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Update User for insert Token
	newUser.Token = token
	newUser, err = u.userService.UpdateUser(newUser)
	if err != nil {
		response := helper.APIResponse("Internal Server Error", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Berhasil Menyimpan Data", "success", http.StatusOK, user.Formatter(newUser, token))
	c.JSON(http.StatusOK, response)
}

// Handler for User Login
func (u *userHandler) LoginUser(c *gin.Context) {
	// Format Input Data & Validation
	var input user.LoginUserInput
	// Input Data From JSON
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Login Failed", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Login Service
	loggedUser, err := u.userService.LoginUser(input)

	if err != nil {
		response := helper.APIResponse("Terjadi Kesalahan", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Generate Token
	token, err := u.authService.GenerateToken(loggedUser.Id)
	if err != nil {
		response := helper.APIResponse("Failed generate token", "error", http.StatusBadRequest, gin.H{"errors": err.Error()})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Barhasil Login", "success", http.StatusOK, user.Formatter(loggedUser, token))
	c.JSON(http.StatusOK, response)
}

// Handler for Check Email Avalilable
func (u *userHandler) CheckEmail(c *gin.Context) {
	// Format Input & Validation
	var input user.EmailInput
	// Input Data from JSON
	err := c.ShouldBindJSON(&input)

	if err != nil {
		response := helper.APIResponse("Email check failed", "error", http.StatusUnprocessableEntity, gin.H{"errors": helper.FormatValidationError(err)})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// Check Email Service
	is_available, err := u.userService.CheckEmail(input)
	if err != nil {
		response := helper.APIResponse("Email check failed", "error", http.StatusUnprocessableEntity, gin.H{"errors": "Internal server error"})
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	message := "Email address has been registered"

	if is_available {
		message = "Email address is available"
	}

	response := helper.APIResponse(message, "success", http.StatusOK, gin.H{"is_available": is_available})
	c.JSON(http.StatusOK, response)
}

// Handler for Upload Image Profile User
func (u *userHandler) UploadImage(c *gin.Context) {
	// Data from form file
	file, err := c.FormFile("photo")

	if err != nil {
		response := helper.APIResponse("Upload Failed", "error", http.StatusBadRequest, gin.H{"is_uploaded": false})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Get Data User Id Session
	id := c.MustGet("currentUser").(user.User).Id
	// Create Path Image
	path := fmt.Sprintf("images/%d-%s", id, file.Filename)
	// Save Data File to Path
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		response := helper.APIResponse("Upload Failed", "error", http.StatusBadRequest, gin.H{"is_uploaded": false})
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// Update Photo Profile User Service
	_, err = u.userService.UploadPhoto(id, path)

	if err != nil {
		response := helper.APIResponse("Upload Failed", "error", http.StatusBadRequest, gin.H{"is_uploaded": false})
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Upload Success", "success", http.StatusOK, gin.H{"is_uploaded": true})
	c.JSON(http.StatusOK, response)
}

// Handler for Get Data Current User
func (u *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)

	response := helper.APIResponse("Upload Success", "success", http.StatusOK, user.Formatter(currentUser, ""))
	c.JSON(http.StatusOK, response)
}

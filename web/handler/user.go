package handler

import (
	"bwa_golang/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(u user.Service) *userHandler {
	return &userHandler{u}
}

func (u *userHandler) Index(c *gin.Context) {
	users, err := u.userService.GetAllUser()

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "index_user.html", gin.H{"users": users})
}
func (u *userHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "create_user.html", nil)
}
func (u *userHandler) Create(c *gin.Context) {
	var inputUser user.CreateUserAdmin

	err := c.ShouldBind(&inputUser)
	if err != nil {
		inputUser.Error = err.Error()
		c.HTML(http.StatusUnprocessableEntity, "create_user.html", inputUser)
		return
	}
	file, err := c.FormFile("photo")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := user.RegisterUserInput{
		Name:       inputUser.Name,
		Email:      inputUser.Email,
		Password:   inputUser.Password,
		Profession: inputUser.Profession,
	}

	dataUser, err := u.userService.RegisterUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	// Upload Photo
	if file != nil {
		// Create Path Image
		path := fmt.Sprintf("images/%d-%s", dataUser.Id, file.Filename)
		// Save Data File to Path
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		// Update Photo Profile User Service
		_, err = u.userService.UploadPhoto(dataUser.Id, path)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
	}
	c.Redirect(http.StatusFound, "/admin_user")
}
func (u *userHandler) Edit(c *gin.Context) {
	// var uri_id user.GetUserIdUri
	id := c.Param("id")
	id_user, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	user, err := u.userService.GetUserById(id_user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.HTML(http.StatusOK, "edit_user.html", gin.H{"user": user})
}
func (u *userHandler) Update(c *gin.Context) {
	var inputUser user.CreateUserAdmin

	err := c.ShouldBind(&inputUser)

	if err != nil {
		inputUser.Error = err.Error()
		// c.HTML(http.StatusUnprocessableEntity, "edit_user.html", nil)
		return
	}

	id := c.Param("id")
	id_user, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	dataUser, err := u.userService.GetUserById(id_user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	file, err := c.FormFile("photo")

	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// Upload Photo
	if file != nil {
		// Create Path Image
		path := fmt.Sprintf("images/%d-%s", dataUser.Id, file.Filename)
		// Save Data File to Path
		err = c.SaveUploadedFile(file, path)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		dataUser.Image = path
	}
	dataUser.Name = inputUser.Name
	dataUser.Email = inputUser.Email
	dataUser.Profession = inputUser.Profession

	_, err = u.userService.UpdateUser(dataUser)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/admin_user")
}
func (u *userHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	id_user, err := strconv.Atoi(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	err = u.userService.DeleteUser(id_user)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/admin_user")
}

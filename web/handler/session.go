package handler

import (
	"bwa_golang/user"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type sessionHandler struct {
	userService user.Service
}

func NewSessionHandler(u user.Service) *sessionHandler {
	return &sessionHandler{u}
}

func (s *sessionHandler) Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func (s *sessionHandler) Sign_in(c *gin.Context) {
	var input user.LoginAdminInput

	err := c.ShouldBind(&input)
	if err != nil {
		c.HTML(http.StatusUnprocessableEntity, "error.html", nil)
		return
	}

	dataUser, err := s.userService.LoginUser(user.LoginUserInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil || dataUser.Role != 1 {
		c.Redirect(http.StatusFound, "/login_admin")
		return
	}

	session := sessions.Default(c)
	session.Set("adminId", dataUser.Id)
	session.Save()
	c.Redirect(http.StatusFound, "/admin_user")
}
func (s *sessionHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/login_admin")
}

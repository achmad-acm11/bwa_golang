package user

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	Profession string `json:"profession" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginAdminInput struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required"`
}

type EmailInput struct {
	Email string `json:"email" binding:"required,email"`
}

type CreateUserAdmin struct {
	Name       string `form:"name" binding:"required"`
	Email      string `form:"email" binding:"required,email"`
	Password   string `form:"password" binding:"required"`
	Profession string `form:"profession" binding:"required"`
	Error      string
}

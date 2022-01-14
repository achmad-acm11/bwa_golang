package user

type UserFormatter struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Profession string `json:"profession"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	ImageURL   string `json:"image_url"`
}

func Formatter(user User, token string) UserFormatter {
	userFormat := UserFormatter{
		Id:         user.Id,
		Name:       user.Name,
		Profession: user.Profession,
		Email:      user.Email,
		Token:      token,
		ImageURL:   user.Image,
	}
	return userFormat
}

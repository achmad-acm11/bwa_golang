package project

import "strings"

type ProjectFormatter struct {
	Id                int    `json:"id"`
	User_id           int    `json:"user_id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Short_description string `json:"short_description"`
	Slug              string `json:"slug"`
	Goal_amount       int    `json:"goal_amount"`
	Current_amount    int    `json:"current_amount"`
	Backer_count      int    `json:"backer_count"`
	Image_url         string `json:"image_url"`
}
type ProjectDetailFormatter struct {
	Id                int              `json:"id"`
	User_id           int              `json:"user_id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Short_description string           `json:"short_description"`
	Perks             []string         `json:"perks"`
	Slug              string           `json:"slug"`
	Goal_amount       int              `json:"goal_amount"`
	Current_amount    int              `json:"current_amount"`
	Backer_count      int              `json:"backer_count"`
	Image_url         string           `json:"image_url"`
	User              UserFormatter    `json:"user"`
	Images            []ImageFormatter `json:"images"`
}
type UserFormatter struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Image_url string `json:"image_url"`
}
type ImageFormatter struct {
	Id         int    `json:"id"`
	Image_url  string `json:"image_url"`
	Is_primary bool   `json:"is_primary"`
}

func FormatterAll(project []Project) []ProjectFormatter {
	var projectFormat []ProjectFormatter
	for _, v := range project {
		projectFormat = append(projectFormat, Formatter(v))
	}
	return projectFormat
}
func FormatterDetail(project Project) ProjectDetailFormatter {
	projectFormat := ProjectDetailFormatter{
		Id:                project.Id,
		User_id:           project.Id_user,
		Name:              project.Project_name,
		Description:       project.Description,
		Short_description: project.Short_description,
		Slug:              project.Slug,
		Goal_amount:       project.Goal_amount,
		Current_amount:    project.Current_amount,
		Backer_count:      project.Backer_count,
	}
	for _, v := range strings.Split(project.Perks, ",") {
		projectFormat.Perks = append(projectFormat.Perks, strings.TrimSpace(v))
	}
	if len(project.Images) > 0 {
		projectFormat.Image_url = project.Images[0].Image
		for _, v := range project.Images {
			primary := false
			if v.Is_primary == 1 {
				primary = true
			}
			imageFormat := ImageFormatter{
				Id:         v.Id,
				Image_url:  v.Image,
				Is_primary: primary,
			}
			projectFormat.Images = append(projectFormat.Images, imageFormat)
		}
	}
	projectFormat.User = UserFormatter{
		Id:        project.User.Id,
		Name:      project.User.Name,
		Image_url: project.User.Image,
	}
	return projectFormat
}
func Formatter(project Project) ProjectFormatter {
	projectFormat := ProjectFormatter{
		Id:                project.Id,
		User_id:           project.Id_user,
		Name:              project.Project_name,
		Description:       project.Description,
		Short_description: project.Short_description,
		Slug:              project.Slug,
		Goal_amount:       project.Goal_amount,
		Current_amount:    project.Current_amount,
		Backer_count:      project.Backer_count,
	}
	if len(project.Images) > 0 {
		projectFormat.Image_url = project.Images[0].Image
	}

	return projectFormat
}

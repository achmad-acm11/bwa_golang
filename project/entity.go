package project

import (
	"bwa_golang/helper"
	"bwa_golang/image"
	"bwa_golang/user"
)

// Project Model
type Project struct {
	Id                int
	Id_user           int
	Project_name      string
	Description       string
	Short_description string
	Perks             string
	Slug              string
	Goal_amount       int
	Current_amount    int
	Backer_count      int
	Images            []image.Image `gorm:"ForeignKey:Id_project"`
	User              user.User     `gorm:"ForeignKey:Id_user"`
}

func (p Project) FormatAmount(number int) string {
	return helper.FormatIDR(number)
}

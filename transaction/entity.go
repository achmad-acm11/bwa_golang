package transaction

import (
	"bwa_golang/project"
	"bwa_golang/user"
	"time"
)

type Transaction struct {
	Id           int
	Id_user      int
	Id_project   int
	Payment_url  string
	Code         string
	Amount       int
	Status       int
	Created_date time.Time
	Updated_date time.Time
	User         user.User       `gorm:"ForeignKey:Id_user"`
	Project      project.Project `gorm:"ForeignKey:Id_project"`
}

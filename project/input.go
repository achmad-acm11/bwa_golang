package project

import "bwa_golang/user"

type CreateProjectInput struct {
	Project_name      string `json:"name" binding:"required"`
	Description       string `json:"description" binding:"required"`
	Short_description string `json:"short_description" binding:"required"`
	Perks             string `json:"perks" binding:"required"`
	Goal_amount       int    `json:"goal_amount" binding:"required"`
	Slug              string
	User              user.User
}

type GetProjectIdUri struct {
	Id int `uri:"id" binding:"required"`
}

type CreateProjectAdmin struct {
	Project_name      string `form:"name" binding:"required"`
	Description       string `form:"description" binding:"required"`
	Short_description string `form:"short_description" binding:"required"`
	Perks             string `form:"perks" binding:"required"`
	Goal_amount       int    `form:"goal_amount" binding:"required"`
	Slug              string
	UserId            int `form:"user_id" binding:"required"`
	Users             []user.User
	Error             string
}

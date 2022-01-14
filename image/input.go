package image

type GetImageIdUri struct {
	Id int `uri:"id" binding:"required"`
}

type CreateImageInput struct {
	Is_primary bool `form:"is_primary"`
	Image      string
}

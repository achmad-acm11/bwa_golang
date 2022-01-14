package image

type ImageFormatter struct {
	Id         int    `json:"id"`
	Image_url  string `json:"image_url"`
	Is_primary bool   `json:"is_primary"`
}

func FormatterImage(image Image) ImageFormatter {
	is_primary := false
	if image.Is_primary == 1 {
		is_primary = true
	}
	imageFormat := ImageFormatter{
		Id:         image.Id,
		Image_url:  image.Image,
		Is_primary: is_primary,
	}
	return imageFormat
}

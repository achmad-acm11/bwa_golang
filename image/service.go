package image

type Service interface {
	CreateImage(input CreateImageInput, id_project int) (Image, error)
}

type service struct {
	imageRepository *repository
}

func NewService(repository *repository) *service {
	return &service{repository}
}

// Create Image Service
func (s *service) CreateImage(input CreateImageInput, id_project int) (Image, error) {
	// Is Clear Primary image all
	is_primary := 0
	if input.Is_primary {
		err := s.imageRepository.ClearPrimaryImage(id_project)
		if err != nil {
			return Image{}, err
		}
		is_primary = 1
	}
	// Mapping Input Image to Image Model
	image := Image{
		Id_project: id_project,
		Image:      input.Image,
		Is_primary: is_primary,
	}
	// Create Image
	image, err := s.imageRepository.CreateImage(image)

	if err != nil {
		return image, err
	}
	return image, nil
}

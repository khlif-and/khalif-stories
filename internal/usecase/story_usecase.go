package usecase

import (
	"errors"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"

	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type storyUseCase struct {
	repo     domain.StoryRepository
	uploader *utils.AzureUploader
}

func NewStoryUseCase(repo domain.StoryRepository, uploader *utils.AzureUploader) domain.StoryUseCase {
	return &storyUseCase{
		repo:     repo,
		uploader: uploader,
	}
}

func (u *storyUseCase) CreateCategory(name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	// 1. Cek Duplikat Nama Dulu
	existing, err := u.repo.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	// 2. Jika Aman, Baru Lanjut Upload & Proses Warna
	var imageURL string
	dominantColor := "#000000"

	if file != nil {
		color, err := utils.ExtractDominantColor(file)
		if err == nil {
			dominantColor = color
		}

		file.Seek(0, 0)

		ext := filepath.Ext(header.Filename)
		filename := uuid.New().String() + ext

		url, err := u.uploader.UploadToContainer(file, "category", filename)
		if err != nil {
			return nil, err
		}
		imageURL = url
	}

	category := &domain.Category{
		Name:          name,
		ImageURL:      imageURL,
		DominantColor: dominantColor,
	}

	err = u.repo.CreateCategory(category)
	return category, err
}

func (u *storyUseCase) GetAllCategories() ([]domain.Category, error) {
	return u.repo.GetAllCategories()
}

func (u *storyUseCase) CreateStory(title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	var thumbURL string
	dominantColor := "#000000"

	if file != nil {
		color, err := utils.ExtractDominantColor(file)
		if err == nil {
			dominantColor = color
		}

		file.Seek(0, 0)

		ext := filepath.Ext(header.Filename)
		filename := "stories/thumbnails/" + uuid.New().String() + ext

		url, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		thumbURL = url
	}

	story := &domain.Story{
		Title:         title,
		Description:   desc,
		CategoryID:    categoryID,
		ThumbnailURL:  thumbURL,
		DominantColor: dominantColor,
		Status:        "Draft",
	}

	if err := u.repo.CreateStory(story); err != nil {
		return nil, err
	}

	return story, nil
}

func (u *storyUseCase) GetAllStories(page, limit int, sort string) ([]domain.Story, error) {
	return u.repo.GetAllStories(page, limit, sort)
}

func (u *storyUseCase) DeleteStory(id uint) error {
	return u.repo.DeleteStory(id)
}

func (u *storyUseCase) AddSlide(storyID uint, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	count, err := u.repo.CountSlides(storyID)
	if err != nil {
		return nil, err
	}

	if count >= 20 {
		return nil, errors.New("maximum slide limit reached (20)")
	}

	var imageURL string
	if file != nil {
		ext := filepath.Ext(header.Filename)
		filename := "stories/slides/" + uuid.New().String() + ext

		url, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		imageURL = url
	}

	slide := &domain.Slide{
		StoryID:  storyID,
		Content:  content,
		Sequence: sequence,
		ImageURL: imageURL,
	}

	if err := u.repo.CreateSlide(slide); err != nil {
		return nil, err
	}

	return slide, nil
}
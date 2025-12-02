package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"khalif-stories/internal/domain"
	"khalif-stories/internal/repository"
	"khalif-stories/pkg/utils"

)

type storyUseCase struct {
	repo       domain.StoryRepository
	redisRepo  *repository.RedisRepo
	searchRepo domain.SearchRepository
	uploader   *utils.AzureUploader
}

func NewStoryUseCase(repo domain.StoryRepository, redisRepo *repository.RedisRepo, searchRepo domain.SearchRepository, uploader *utils.AzureUploader) domain.StoryUseCase {
	return &storyUseCase{
		repo:       repo,
		redisRepo:  redisRepo,
		searchRepo: searchRepo,
		uploader:   uploader,
	}
}

func (u *storyUseCase) SearchCategories(query string) ([]domain.Category, error) {
	return u.searchRepo.SearchCategories(query)
}

func (u *storyUseCase) CreateCategory(name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	existing, err := u.repo.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category with this name already exists")
	}

	var imageURL string
	dominantColor := "#000000"

	fileBytes, err := utils.ReadMultipartFileToBytes(file)
	if err != nil {
		return nil, err
	}

	if fileBytes != nil {
		ext := filepath.Ext(header.Filename)
		filename := uuid.New().String() + ext

		url, err := u.uploader.UploadToContainer(file, "category", filename)
		if err != nil {
			return nil, err
		}
		imageURL = url

		if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
			dominantColor = color
		}
	}

	category := &domain.Category{
		UUID:          uuid.New().String(),
		Name:          name,
		ImageURL:      imageURL,
		DominantColor: dominantColor,
	}

	err = u.repo.CreateCategory(category)
	if err != nil {
		return nil, err
	}

	go u.searchRepo.IndexCategory(category)

	return category, nil
}

func (u *storyUseCase) UpdateCategory(uuidStr string, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	category, err := u.repo.GetCategoryByUUID(uuidStr)
	if err != nil {
		return nil, err
	}

	if name != "" && name != category.Name {
		existing, err := u.repo.GetCategoryByName(name)
		if err == nil && existing != nil && existing.ID != category.ID {
			return nil, errors.New("category name already used")
		}
		category.Name = name
	}

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}
		if fileBytes != nil {
			ext := filepath.Ext(header.Filename)
			filename := uuid.New().String() + ext
			url, err := u.uploader.UploadToContainer(file, "category", filename)
			if err != nil {
				return nil, err
			}
			category.ImageURL = url
			if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
				category.DominantColor = color
			}
		}
	}

	if err := u.repo.UpdateCategory(category); err != nil {
		return nil, err
	}

	go u.searchRepo.IndexCategory(category)

	return category, nil
}

func (u *storyUseCase) DeleteCategory(uuidStr string) error {
	err := u.repo.DeleteCategory(uuidStr)
	if err != nil {
		return err
	}

	go u.searchRepo.DeleteCategoryIndex(uuidStr)

	return nil
}

func (u *storyUseCase) GetAllCategories() ([]domain.Category, error) {
	return u.repo.GetAllCategories()
}

func (u *storyUseCase) GetCategory(uuidStr string) (*domain.Category, error) {
	return u.repo.GetCategoryByUUID(uuidStr)
}

func (u *storyUseCase) CreateStory(title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	var thumbURL string
	dominantColor := "#000000"

	fileBytes, err := utils.ReadMultipartFileToBytes(file)
	if err != nil {
		return nil, err
	}

	if fileBytes != nil {
		ext := filepath.Ext(header.Filename)
		filename := "stories/thumbnails/" + uuid.New().String() + ext

		url, err := u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
		thumbURL = url

		if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
			dominantColor = color
		}
	}

	story := &domain.Story{
		UUID:          uuid.New().String(),
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
	cacheKey := fmt.Sprintf("stories:p%d:l%d:s%s", page, limit, sort)

	cachedData, err := u.redisRepo.Get(cacheKey)
	if err == nil && cachedData != "" {
		var stories []domain.Story
		if err := json.Unmarshal([]byte(cachedData), &stories); err == nil {
			return stories, nil
		}
	}

	stories, err := u.repo.GetAllStories(page, limit, sort)
	if err != nil {
		return nil, err
	}

	if jsonData, err := json.Marshal(stories); err == nil {
		u.redisRepo.Set(cacheKey, jsonData, 30*time.Minute)
	}

	return stories, nil
}

func (u *storyUseCase) DeleteStory(uuid string) error {
	return u.repo.DeleteStory(uuid)
}

func (u *storyUseCase) AddSlide(storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	story, err := u.repo.GetStoryByUUID(storyUUID)
	if err != nil {
		return nil, err
	}

	count, err := u.repo.CountSlides(story.ID)
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
		StoryID:  story.ID,
		Content:  content,
		Sequence: sequence,
		ImageURL: imageURL,
	}

	if err := u.repo.CreateSlide(slide); err != nil {
		return nil, err
	}

	return slide, nil
}
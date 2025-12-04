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

type StoryUC struct {
	repo      domain.StoryRepository
	redisRepo *repository.RedisRepo
	uploader  *utils.AzureUploader
}

func NewStoryUseCase(repo domain.StoryRepository, redisRepo *repository.RedisRepo, uploader *utils.AzureUploader) domain.StoryUseCase {
	return &StoryUC{repo: repo, redisRepo: redisRepo, uploader: uploader}
}

func (u *StoryUC) refreshCache() {
	u.redisRepo.DeletePrefix("stories:")

	defaultPage := 1
	defaultLimit := 10
	defaultSort := "created_at desc"

	stories, err := u.repo.GetAll(defaultPage, defaultLimit, defaultSort)
	if err != nil {
		return
	}

	if data, err := json.Marshal(stories); err == nil {
		cacheKey := fmt.Sprintf("stories:p%d:l%d:s%s", defaultPage, defaultLimit, defaultSort)
		u.redisRepo.Set(cacheKey, data, 30*time.Minute)
	}
}

func (u *StoryUC) Create(title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	var thumbURL string
	dominantColor := "#000000"

	fileBytes, err := utils.ReadMultipartFileToBytes(file)
	if err != nil {
		return nil, err
	}
	if fileBytes != nil {
		filename := "stories/thumbnails/" + uuid.New().String() + filepath.Ext(header.Filename)
		thumbURL, err = u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
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

	if err := u.repo.Create(story); err != nil {
		return nil, err
	}

	go u.refreshCache()
	
	return story, nil
}

func (u *StoryUC) GetAll(page, limit int, sort string) ([]domain.Story, error) {
	cacheKey := fmt.Sprintf("stories:p%d:l%d:s%s", page, limit, sort)
	if cached, err := u.redisRepo.Get(cacheKey); err == nil && cached != "" {
		var stories []domain.Story
		if json.Unmarshal([]byte(cached), &stories) == nil {
			return stories, nil
		}
	}

	stories, err := u.repo.GetAll(page, limit, sort)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(stories); err == nil {
		u.redisRepo.Set(cacheKey, data, 30*time.Minute)
	}
	return stories, nil
}

func (u *StoryUC) Search(query string) (*[]domain.Story, error) {
	stories, err := u.repo.Search(query)
	return &stories, err
}

func (u *StoryUC) Delete(uuid string) error {
	if err := u.repo.Delete(uuid); err != nil {
		return err
	}
	
	go u.refreshCache()
	
	return nil
}

func (u *StoryUC) AddSlide(storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	story, err := u.repo.GetByUUID(storyUUID)
	if err != nil {
		return nil, err
	}

	count, err := u.repo.CountSlides(story.ID)
	if err != nil {
		return nil, err
	}
	if count >= 20 {
		return nil, errors.New("slide limit reached")
	}

	var imageURL string
	if file != nil {
		filename := "stories/slides/" + uuid.New().String() + filepath.Ext(header.Filename)
		imageURL, err = u.uploader.UploadFile(file, filename)
		if err != nil {
			return nil, err
		}
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
	
	go u.refreshCache()

	return slide, nil
}
package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/repository"
	"khalif-stories/pkg/utils"

)

type StoryUC struct {
	cfg       *config.Config
	repo      domain.StoryRepository
	redisRepo *repository.RedisRepo
	uploader  *utils.AzureUploader
}

func NewStoryUseCase(cfg *config.Config, repo domain.StoryRepository, redisRepo *repository.RedisRepo, uploader *utils.AzureUploader) domain.StoryUseCase {
	return &StoryUC{cfg: cfg, repo: repo, redisRepo: redisRepo, uploader: uploader}
}

func (u *StoryUC) Create(ctx context.Context, title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	var thumbURL string
	dominantColor := "#000000"

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}
		if fileBytes != nil {
			filename := "stories/thumbnails/" + uuid.New().String() + filepath.Ext(header.Filename)
			thumbURL, err = u.uploader.UploadFile(ctx, file, filename)
			if err != nil {
				return nil, err
			}
			if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
				dominantColor = color
			}
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

	if err := u.repo.Create(ctx, story); err != nil {
		return nil, err
	}

	return story, nil
}

func (u *StoryUC) GetAll(ctx context.Context, page, limit int, sort string) ([]domain.Story, error) {
	cacheKey := fmt.Sprintf("stories:p%d:l%d:s%s", page, limit, sort)
	if cached, err := u.redisRepo.Get(ctx, cacheKey); err == nil && cached != "" {
		var stories []domain.Story
		if json.Unmarshal([]byte(cached), &stories) == nil {
			return stories, nil
		}
	}

	stories, err := u.repo.GetAll(ctx, page, limit, sort)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(stories); err == nil {
		u.redisRepo.Set(ctx, cacheKey, data, 5*time.Minute)
	}
	return stories, nil
}

func (u *StoryUC) Search(ctx context.Context, query string) (*[]domain.Story, error) {
	stories, err := u.repo.Search(ctx, query)
	return &stories, err
}

func (u *StoryUC) Delete(ctx context.Context, uuid string) error {
	return u.repo.Delete(ctx, uuid)
}

func (u *StoryUC) AddSlide(ctx context.Context, storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	story, err := u.repo.GetByUUID(ctx, storyUUID)
	if err != nil {
		return nil, err
	}

	count, err := u.repo.CountSlides(ctx, story.ID)
	if err != nil {
		return nil, err
	}
	
	if count >= int64(u.cfg.SlideLimit) {
		return nil, errors.New("slide limit reached")
	}

	var imageURL string
	if file != nil {
		filename := "stories/slides/" + uuid.New().String() + filepath.Ext(header.Filename)
		imageURL, err = u.uploader.UploadFile(ctx, file, filename)
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

	if err := u.repo.CreateSlide(ctx, slide); err != nil {
		return nil, err
	}

	return slide, nil
}
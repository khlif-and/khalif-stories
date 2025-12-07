package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/internal/repository"
	"khalif-stories/pkg/utils"

)

type StoryUC struct {
	cfg          *config.Config
	repo         domain.StoryRepository
	categoryRepo domain.CategoryRepository
	redisRepo    *repository.RedisRepo
	uploader     *utils.AzureUploader
}

func NewStoryUseCase(cfg *config.Config, repo domain.StoryRepository, categoryRepo domain.CategoryRepository, redisRepo *repository.RedisRepo, uploader *utils.AzureUploader) domain.StoryUseCase {
	return &StoryUC{cfg: cfg, repo: repo, categoryRepo: categoryRepo, redisRepo: redisRepo, uploader: uploader}
}

func (u *StoryUC) Create(ctx context.Context, title, desc string, categoryUUID string, userID string, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	if isDup, _ := u.repo.CheckDuplicate(ctx, title, desc); isDup {
		return nil, errors.New("story with same title already exists")
	}

	cat, err := u.categoryRepo.GetByUUID(ctx, categoryUUID)
	if err != nil || cat == nil {
		return nil, errors.New("invalid category")
	}

	story := &domain.Story{
		UUID:        uuid.New().String(),
		Title:       title,
		Description: desc,
		CategoryID:  cat.ID,
		Category:    *cat,
		UserID:      userID,
		Status:      "Pending_Upload",
	}

	if err := u.repo.Create(ctx, story); err != nil {
		return nil, err
	}

	thumbURL, domColor, err := utils.UploadAndAnalyzeImage(ctx, u.uploader, file, header, u.cfg.AzureContainerStoriesName, u.cfg.StoriesThumbPath, story.UUID)
	if err != nil {
		u.repo.Delete(ctx, story.UUID)
		return nil, err
	}

	story.ThumbnailURL = thumbURL
	story.DominantColor = domColor
	story.Status = domain.StatusDraft

	if err := u.repo.Update(ctx, story); err != nil {
		u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerStoriesName, thumbURL)
		u.repo.Delete(ctx, story.UUID)
		return nil, err
	}

	if u.redisRepo != nil {
		_ = u.redisRepo.DeletePrefix(ctx, domain.CacheKeyStoryPrefix)
	}

	return story, nil
}

func (u *StoryUC) Update(ctx context.Context, storyUUID string, title, desc, categoryUUID, status string, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	story, err := u.repo.GetByUUID(ctx, storyUUID)
	if err != nil || story == nil {
		return nil, errors.New("story not found")
	}

	oldThumbURL := story.ThumbnailURL

	if title != "" {
		story.Title = title
	}
	if desc != "" {
		story.Description = desc
	}
	if status != "" {
		story.Status = status
	}

	if categoryUUID != "" {
		if cat, _ := u.categoryRepo.GetByUUID(ctx, categoryUUID); cat != nil {
			story.CategoryID = cat.ID
			story.Category = *cat
		}
	}

	newThumbURL, newColor, err := utils.UploadAndAnalyzeImage(ctx, u.uploader, file, header, u.cfg.AzureContainerStoriesName, u.cfg.StoriesThumbPath, uuid.New().String())
	if err != nil {
		return nil, err
	}

	if newThumbURL != "" {
		story.ThumbnailURL = newThumbURL
		story.DominantColor = newColor
	}

	story.UpdatedAt = time.Now()

	if err := u.repo.Update(ctx, story); err != nil {
		if newThumbURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerStoriesName, newThumbURL)
		}
		return nil, err
	}

	if newThumbURL != "" && oldThumbURL != "" {
		u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerStoriesName, oldThumbURL)
	}

	if u.redisRepo != nil {
		_ = u.redisRepo.DeletePrefix(ctx, domain.CacheKeyStoryPrefix)
	}

	return story, nil
}

func (u *StoryUC) Delete(ctx context.Context, uuid string) error {
	story, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if story == nil {
		return nil
	}

	if err := u.repo.Delete(ctx, uuid); err != nil {
		return err
	}

	if story.ThumbnailURL != "" {
		u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerStoriesName, story.ThumbnailURL)
	}
	for _, slide := range story.Slides {
		if slide.ImageURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainer, slide.ImageURL)
		}
	}

	if u.redisRepo != nil {
		_ = u.redisRepo.DeletePrefix(ctx, domain.CacheKeyStoryPrefix)
	}

	return nil
}

func (u *StoryUC) AddSlide(ctx context.Context, storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*domain.Slide, error) {
	story, err := u.repo.GetByUUID(ctx, storyUUID)
	if err != nil {
		return nil, err
	}

	count, _ := u.repo.CountSlides(ctx, story.ID)
	if count >= int64(u.cfg.SlideLimit) {
		return nil, errors.New("slide limit reached")
	}

	imageURL, _, err := utils.UploadAndAnalyzeImage(ctx, u.uploader, file, header, u.cfg.AzureContainer, u.cfg.StoriesSlidePath, uuid.New().String())
	if err != nil {
		return nil, err
	}

	slide := &domain.Slide{
		StoryID:  &story.ID,
		Content:  content,
		Sequence: sequence,
		ImageURL: imageURL,
	}

	if err := u.repo.CreateSlide(ctx, slide); err != nil {
		if imageURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainer, imageURL)
		}
		return nil, err
	}

	return slide, nil
}

func (u *StoryUC) GetAll(ctx context.Context, page, limit int, sort string) ([]domain.Story, error) {
	cacheKey := fmt.Sprintf("stories:p%d:l%d:s%s", page, limit, sort)
	if cached, _ := u.redisRepo.Get(ctx, cacheKey); cached != "" {
		var stories []domain.Story
		if json.Unmarshal([]byte(cached), &stories) == nil {
			return stories, nil
		}
	}

	stories, err := u.repo.GetAll(ctx, page, limit, sort)
	if err == nil {
		if data, err := json.Marshal(stories); err == nil {
			u.redisRepo.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}
	return stories, err
}

func (u *StoryUC) GetByUUID(ctx context.Context, uuid string) (*domain.Story, error) {
	story, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, domain.ErrNotFound
	}
	return story, nil
}

func (u *StoryUC) Search(ctx context.Context, query string) ([]domain.Story, error) {
	return u.repo.Search(ctx, query)
}

// --- [INI FUNGSI YANG KITA TAMBAHKAN] ---
func (u *StoryUC) GetRecommendations(ctx context.Context, userID string) ([]domain.Recommendation, error) {
	return u.repo.GetRecommendations(ctx, userID)
}
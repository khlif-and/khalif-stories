package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
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
	isDuplicate, err := u.repo.CheckDuplicate(ctx, title, desc)
	if err != nil {
		return nil, err
	}
	if isDuplicate {
		return nil, errors.New("story with the same title already exists")
	}

	category, err := u.categoryRepo.GetByUUID(ctx, categoryUUID)
	if err != nil {
		return nil, errors.New("invalid category or database error")
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	var thumbURL string
	dominantColor := "#000000"
	
	var uploadedFilename string
	var azClient *azblob.Client
	var azContainer string

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}
		if fileBytes != nil {
			folderPath := u.cfg.StoriesThumbPath
			uploadedFilename = folderPath + uuid.New().String() + filepath.Ext(header.Filename)

			connectionString := u.cfg.AzureConnStr
			azContainer = u.cfg.AzureContainerStoriesName

			if connectionString == "" || azContainer == "" {
				return nil, errors.New("azure configuration is missing")
			}

			azClient, err = azblob.NewClientFromConnectionString(connectionString, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create azure client: %w", err)
			}

			_, err = azClient.UploadBuffer(ctx, azContainer, uploadedFilename, fileBytes, &azblob.UploadBufferOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to upload to azure: %w", err)
			}

			thumbURL = fmt.Sprintf("%s/%s/%s", azClient.URL(), azContainer, uploadedFilename)

			if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
				dominantColor = color
			}
		}
	}

	story := &domain.Story{
		UUID:          uuid.New().String(),
		Title:         title,
		Description:   desc,
		CategoryID:    category.ID,
		Category:      *category,
		UserID:        userID,
		ThumbnailURL:  thumbURL,
		DominantColor: dominantColor,
		Status:        "Draft",
	}

	if err := u.repo.Create(ctx, story); err != nil {
		if uploadedFilename != "" && azClient != nil {
			azClient.DeleteBlob(ctx, azContainer, uploadedFilename, nil)
		}
		return nil, err
	}

	return story, nil
}

func (u *StoryUC) Update(ctx context.Context, storyUUID string, title, desc, categoryUUID, status string, file multipart.File, header *multipart.FileHeader) (*domain.Story, error) {
	story, err := u.repo.GetByUUID(ctx, storyUUID)
	if err != nil {
		return nil, err
	}
	if story == nil {
		return nil, errors.New("story not found")
	}

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
		category, err := u.categoryRepo.GetByUUID(ctx, categoryUUID)
		if err != nil {
			return nil, errors.New("invalid category or database error")
		}
		if category != nil {
			story.CategoryID = category.ID
			story.Category = *category
		}
	}

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}

		if fileBytes != nil {
			folderPath := u.cfg.StoriesThumbPath
			uploadedFilename := folderPath + uuid.New().String() + filepath.Ext(header.Filename)

			connectionString := u.cfg.AzureConnStr
			azContainer := u.cfg.AzureContainerStoriesName

			if connectionString == "" || azContainer == "" {
				return nil, errors.New("azure configuration is missing")
			}

			azClient, err := azblob.NewClientFromConnectionString(connectionString, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to create azure client: %w", err)
			}

			_, err = azClient.UploadBuffer(ctx, azContainer, uploadedFilename, fileBytes, &azblob.UploadBufferOptions{})
			if err != nil {
				return nil, fmt.Errorf("failed to upload to azure: %w", err)
			}

			newThumbURL := fmt.Sprintf("%s/%s/%s", azClient.URL(), azContainer, uploadedFilename)
			
			if story.ThumbnailURL != "" {
				oldBlobName := extractBlobNameFromURL(story.ThumbnailURL, azContainer)
				if oldBlobName != "" {
					azClient.DeleteBlob(ctx, azContainer, oldBlobName, nil)
				}
			}

			story.ThumbnailURL = newThumbURL

			if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
				story.DominantColor = color
			}
		}
	}

	story.UpdatedAt = time.Now()
	
	if err := u.repo.Update(ctx, story); err != nil {
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
		folderPath := u.cfg.StoriesSlidePath
		filename := folderPath + uuid.New().String() + filepath.Ext(header.Filename)

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

func extractBlobNameFromURL(fullURL, containerName string) string {
	parts := strings.Split(fullURL, containerName+"/")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
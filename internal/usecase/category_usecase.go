package usecase

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"time"

	"github.com/google/uuid"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type CategoryUC struct {
	categoryRepo domain.CategoryRepository
	redisRepo    domain.RedisRepository
	uploader     *utils.AzureUploader
	cfg          *config.Config
}

func NewCategoryUseCase(repo domain.CategoryRepository, redis domain.RedisRepository, uploader *utils.AzureUploader) *CategoryUC {
	cfg := config.LoadConfig() 
	
	return &CategoryUC{
		categoryRepo: repo,
		redisRepo:    redis,
		uploader:     uploader,
		cfg:          cfg,
	}
}

func (uc *CategoryUC) Create(ctx context.Context, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	existing, _ := uc.categoryRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, domain.ErrConflict
	}

	category := &domain.Category{
		UUID: uuid.New().String(),
		Name: name,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	if file != nil {
		imageURL, domColor, err := utils.UploadAndAnalyzeImage(ctx, uc.uploader, file, header, uc.cfg.AzureContainer, "categories/", category.UUID)
		
		if err != nil {
			uc.categoryRepo.Delete(ctx, category.UUID)
			return nil, err
		}

		category.ImageURL = imageURL
		category.DominantColor = domColor

		if err := uc.categoryRepo.Update(ctx, category); err != nil {
			uc.uploader.DeleteFromContainer(ctx, uc.cfg.AzureContainer, imageURL)
			uc.categoryRepo.Delete(ctx, category.UUID)
			return nil, err
		}
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return category, nil
}

func (uc *CategoryUC) Update(ctx context.Context, uuid string, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByUUID(ctx, uuid)
	if err != nil { return nil, err }
	if category == nil { return nil, domain.ErrNotFound }

	oldImageURL := category.ImageURL

	if name != "" && name != category.Name {
		existing, _ := uc.categoryRepo.GetByName(ctx, name)
		if existing != nil && existing.UUID != category.UUID {
			return nil, domain.ErrConflict
		}
		category.Name = name
	}

	var newImageURL string
	if file != nil {
		newUUID := uuid
		
		url, color, err := utils.UploadAndAnalyzeImage(ctx, uc.uploader, file, header, uc.cfg.AzureContainer, "categories/", newUUID)
		if err != nil {
			return nil, err
		}
		
		newImageURL = url
		category.ImageURL = newImageURL
		category.DominantColor = color
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		if newImageURL != "" {
			uc.uploader.DeleteFromContainer(ctx, uc.cfg.AzureContainer, newImageURL)
		}
		return nil, err
	}

	if newImageURL != "" && oldImageURL != "" {
		uc.uploader.DeleteFromContainer(ctx, uc.cfg.AzureContainer, oldImageURL)
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return category, nil
}

func (uc *CategoryUC) Delete(ctx context.Context, uuid string) error {
	category, err := uc.categoryRepo.GetByUUID(ctx, uuid)
	if err != nil { return err }
	if category == nil { return domain.ErrNotFound }

	if err := uc.categoryRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	if category.ImageURL != "" && uc.uploader != nil {
		_ = uc.uploader.DeleteFromContainer(ctx, uc.cfg.AzureContainer, category.ImageURL)
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return nil
}

func (uc *CategoryUC) GetAll(ctx context.Context) ([]domain.Category, error) {
	cacheKey := domain.CacheKeyCategoryAll

	if uc.redisRepo != nil {
		cachedData, err := uc.redisRepo.Get(ctx, cacheKey)
		if err == nil && cachedData != "" {
			var categories []domain.Category
			if err := json.Unmarshal([]byte(cachedData), &categories); err == nil {
				return categories, nil
			}
		}
	}

	categories, err := uc.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if uc.redisRepo != nil {
		if data, err := json.Marshal(categories); err == nil {
			_ = uc.redisRepo.Set(ctx, cacheKey, data, 30*time.Minute)
		}
	}

	return categories, nil
}

func (uc *CategoryUC) Get(ctx context.Context, uuid string) (*domain.Category, error) {
	cat, err := uc.categoryRepo.GetByUUID(ctx, uuid)
	if err != nil { return nil, err }
	if cat == nil { return nil, domain.ErrNotFound }
	return cat, nil
}

func (uc *CategoryUC) Search(ctx context.Context, query string) ([]domain.Category, error) {
	return uc.categoryRepo.Search(ctx, query)
}
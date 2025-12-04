package usecase

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"time"

	"khalif-stories/internal/domain"

)

type CategoryUC struct {
	categoryRepo domain.CategoryRepository
	redisRepo    domain.RedisRepository
	storage      domain.StorageRepository
}

func NewCategoryUseCase(repo domain.CategoryRepository, redis domain.RedisRepository, storage domain.StorageRepository) *CategoryUC {
	return &CategoryUC{
		categoryRepo: repo,
		redisRepo:    redis,
		storage:      storage,
	}
}

func (uc *CategoryUC) Create(ctx context.Context, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	existing, _ := uc.categoryRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, domain.ErrConflict
	}

	var imageURL string
	var err error

	if file != nil && uc.storage != nil {
		imageURL, err = uc.storage.Upload(file, header)
		if err != nil {
			return nil, err
		}
	}

	category := &domain.Category{
		Name:     name,
		ImageURL: imageURL,
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return category, nil
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
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, domain.ErrNotFound
	}
	return cat, nil
}

func (uc *CategoryUC) Search(ctx context.Context, query string) ([]domain.Category, error) {
	return uc.categoryRepo.Search(ctx, query)
}

func (uc *CategoryUC) Update(ctx context.Context, uuid string, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, domain.ErrNotFound
	}

	if name != "" && name != category.Name {
		existing, _ := uc.categoryRepo.GetByName(ctx, name)
		if existing != nil && existing.UUID != category.UUID {
			return nil, domain.ErrConflict
		}
		category.Name = name
	}

	if file != nil && uc.storage != nil {
		if category.ImageURL != "" {
			_ = uc.storage.Delete(category.ImageURL)
		}

		newImageURL, err := uc.storage.Upload(file, header)
		if err != nil {
			return nil, err
		}
		category.ImageURL = newImageURL
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, err
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return category, nil
}

func (uc *CategoryUC) Delete(ctx context.Context, uuid string) error {
	category, err := uc.categoryRepo.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}
	if category == nil {
		return domain.ErrNotFound
	}

	if category.ImageURL != "" && uc.storage != nil {
		_ = uc.storage.Delete(category.ImageURL)
	}

	if err := uc.categoryRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	if uc.redisRepo != nil {
		_ = uc.redisRepo.DeletePrefix(ctx, domain.CacheKeyCategoryAll)
	}

	return nil
}
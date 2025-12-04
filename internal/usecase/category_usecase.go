package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"khalif-stories/internal/domain"
	"khalif-stories/internal/repository"
	"khalif-stories/pkg/utils"

)

type CategoryUC struct {
	repo      domain.CategoryRepository
	redisRepo *repository.RedisRepo
	uploader  *utils.AzureUploader
}

func NewCategoryUseCase(repo domain.CategoryRepository, redisRepo *repository.RedisRepo, uploader *utils.AzureUploader) domain.CategoryUseCase {
	return &CategoryUC{repo: repo, redisRepo: redisRepo, uploader: uploader}
}

func (u *CategoryUC) refreshCache() {
	u.redisRepo.DeletePrefix("categories:")

	cats, err := u.repo.GetAll()
	if err != nil {
		return
	}

	if data, err := json.Marshal(cats); err == nil {
		u.redisRepo.Set("categories:all", data, 1*time.Hour)
	}
}

func (u *CategoryUC) Create(name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	existing, err := u.repo.GetByName(name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("category name exists")
	}

	var imageURL string
	dominantColor := "#000000"

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}
		if fileBytes != nil {
			filename := uuid.New().String() + filepath.Ext(header.Filename)
			imageURL, err = u.uploader.UploadToContainer(file, "category", filename)
			if err != nil {
				return nil, err
			}
			if color, err := utils.ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
				dominantColor = color
			}
		}
	}

	category := &domain.Category{
		UUID:          uuid.New().String(),
		Name:          name,
		ImageURL:      imageURL,
		DominantColor: dominantColor,
	}

	if err := u.repo.Create(category); err != nil {
		return nil, err
	}

	go u.refreshCache()
	
	return category, nil
}

func (u *CategoryUC) Update(uuidStr string, name string, file multipart.File, header *multipart.FileHeader) (*domain.Category, error) {
	category, err := u.repo.GetByUUID(uuidStr)
	if err != nil {
		return nil, err
	}

	if name != "" && name != category.Name {
		existing, err := u.repo.GetByName(name)
		if err == nil && existing != nil && existing.ID != category.ID {
			return nil, errors.New("category name used")
		}
		category.Name = name
	}

	if file != nil {
		fileBytes, err := utils.ReadMultipartFileToBytes(file)
		if err != nil {
			return nil, err
		}
		if fileBytes != nil {
			filename := uuid.New().String() + filepath.Ext(header.Filename)
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

	if err := u.repo.Update(category); err != nil {
		return nil, err
	}

	go u.refreshCache()

	return category, nil
}

func (u *CategoryUC) Delete(uuidStr string) error {
	err := u.repo.Delete(uuidStr)
	if err == nil {
		go u.refreshCache()
		go u.redisRepo.DeletePrefix("stories:") 
	}
	return err
}

func (u *CategoryUC) GetAll() ([]domain.Category, error) {
	cacheKey := "categories:all"
	cached, err := u.redisRepo.Get(cacheKey)
	if err == nil && cached != "" {
		var cats []domain.Category
		if json.Unmarshal([]byte(cached), &cats) == nil {
			return cats, nil
		}
	}

	cats, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(cats); err == nil {
		u.redisRepo.Set(cacheKey, data, 1*time.Hour)
	}
	return cats, nil
}

func (u *CategoryUC) Get(uuidStr string) (*domain.Category, error) {
	return u.repo.GetByUUID(uuidStr)
}

func (u *CategoryUC) Search(query string) ([]domain.Category, error) {
	return u.repo.Search(query)
}
package repository

import (
	"errors"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type CategoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

func (r *CategoryRepo) GetByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Model(&domain.Category{}).Where("name = ?", name).First(&category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) GetByUUID(uuid string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Model(&domain.Category{}).Where("uuid = ?", uuid).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepo) Create(c *domain.Category) error {
	return r.db.Create(c).Error
}

func (r *CategoryRepo) Update(c *domain.Category) error {
	return r.db.Save(c).Error
}

func (r *CategoryRepo) Delete(uuid string) error {
	var category domain.Category
	if err := r.db.Select("id").Where("uuid = ?", uuid).First(&category).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM slides WHERE story_id IN (SELECT id FROM stories WHERE category_id = ?)", category.ID).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ?", category.ID).Delete(&domain.Story{}).Error; err != nil {
			return err
		}
		return tx.Delete(&category).Error
	})
}

func (r *CategoryRepo) GetAll() ([]domain.Category, error) {
	var cats []domain.Category
	err := r.db.Order("id ASC").Find(&cats).Error
	return cats, err
}

func (r *CategoryRepo) Search(query string) ([]domain.Category, error) {
	var categories []domain.Category
	pattern := "%" + query + "%"
	err := r.db.Where("name ILIKE ?", pattern).Limit(10).Find(&categories).Error
	return categories, err
}

func (r *CategoryRepo) UpdateColor(id uint, color string) error {
	return r.db.Model(&domain.Category{}).Where("id = ?", id).Update("dominant_color", color).Error
}
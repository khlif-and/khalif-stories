package repository

import (
	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type StoryRepo struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) *StoryRepo {
	return &StoryRepo{db: db}
}

func (r *StoryRepo) Create(s *domain.Story) error {
	return r.db.Create(s).Error
}

func (r *StoryRepo) GetAll(page, limit int, sort string) ([]domain.Story, error) {
	var stories []domain.Story
	offset := (page - 1) * limit
	err := r.db.Preload("Category").
		Order(sort).Limit(limit).Offset(offset).Find(&stories).Error
	return stories, err
}

func (r *StoryRepo) Search(query string) ([]domain.Story, error) {
	var stories []domain.Story
	pattern := "%" + query + "%"
	err := r.db.Preload("Category").
		Where("title ILIKE ? OR description ILIKE ?", pattern, pattern).
		Limit(20).Find(&stories).Error
	return stories, err
}

func (r *StoryRepo) GetByID(id uint) (*domain.Story, error) {
	var story domain.Story
	err := r.db.First(&story, id).Error
	return &story, err
}

func (r *StoryRepo) GetByUUID(uuid string) (*domain.Story, error) {
	var story domain.Story
	err := r.db.Preload("Slides", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Where("uuid = ?", uuid).First(&story).Error
	return &story, err
}

func (r *StoryRepo) Update(s *domain.Story) error {
	return r.db.Save(s).Error
}

func (r *StoryRepo) UpdateColor(id uint, color string) error {
	return r.db.Model(&domain.Story{}).Where("id = ?", id).Update("dominant_color", color).Error
}

func (r *StoryRepo) Delete(uuid string) error {
	var story domain.Story
	if err := r.db.Select("id").Where("uuid = ?", uuid).First(&story).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("story_id = ?", story.ID).Delete(&domain.Slide{}).Error; err != nil {
			return err
		}
		return tx.Delete(&story).Error
	})
}

func (r *StoryRepo) CreateSlide(s *domain.Slide) error {
	return r.db.Exec("CALL add_slide_safe(?, ?, ?, ?)", s.StoryID, s.ImageURL, s.Content, s.Sequence).Error
}

func (r *StoryRepo) CountSlides(storyID uint) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Story{}).Select("slide_count").Where("id = ?", storyID).Scan(&count).Error
	return count, err
}
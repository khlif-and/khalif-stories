package repository

import (
	"context"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type StoryRepo struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) *StoryRepo {
	return &StoryRepo{db: db}
}

func (r *StoryRepo) Create(ctx context.Context, s *domain.Story) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *StoryRepo) CheckDuplicate(ctx context.Context, title, description string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Story{}).
		Where("title = ?", title).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *StoryRepo) GetAll(ctx context.Context, page, limit int, sort string) ([]domain.Story, error) {
	var stories []domain.Story
	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("Category").
		Order(sort).Limit(limit).Offset(offset).Find(&stories).Error
	return stories, err
}

func (r *StoryRepo) Search(ctx context.Context, query string) ([]domain.Story, error) {
	var stories []domain.Story
	pattern := "%" + query + "%"
	err := r.db.WithContext(ctx).Preload("Category").
		Where("title ILIKE ? OR description ILIKE ?", pattern, pattern).
		Limit(20).Find(&stories).Error
	return stories, err
}

func (r *StoryRepo) GetByID(ctx context.Context, id uint) (*domain.Story, error) {
	var story domain.Story
	err := r.db.WithContext(ctx).First(&story, id).Error
	return &story, err
}

func (r *StoryRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Story, error) {
	var story domain.Story
	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Where("uuid = ?", uuid).
		First(&story).Error
	
	return &story, err
}

func (r *StoryRepo) Update(ctx context.Context, s *domain.Story) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *StoryRepo) UpdateColor(ctx context.Context, id uint, color string) error {
	return r.db.WithContext(ctx).Model(&domain.Story{}).Where("id = ?", id).Update("dominant_color", color).Error
}

func (r *StoryRepo) Delete(ctx context.Context, uuid string) error {
	var story domain.Story
	if err := r.db.WithContext(ctx).Select("id").Where("uuid = ?", uuid).First(&story).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("story_id = ?", story.ID).Delete(&domain.Slide{}).Error; err != nil {
			return err
		}
		return tx.Delete(&story).Error
	})
}

func (r *StoryRepo) CreateSlide(ctx context.Context, s *domain.Slide) error {
	return r.db.WithContext(ctx).Exec("CALL add_slide_safe(?, ?, ?, ?)", s.StoryID, s.ImageURL, s.Content, s.Sequence).Error
}

func (r *StoryRepo) CountSlides(ctx context.Context, storyID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Story{}).Select("slide_count").Where("id = ?", storyID).Scan(&count).Error
	return count, err
}

func (r *StoryRepo) GetRecommendations(ctx context.Context, userID string) ([]domain.Recommendation, error) {
	var recs []domain.Recommendation

	// Menggunakan Preload untuk memuat data Story dan Category terkait
	err := r.db.WithContext(ctx).
		Preload("Story").
		Preload("Story.Category").
		Where("user_id = ?", userID).
		Order("score DESC").
		Limit(10).
		Find(&recs).Error

	return recs, err
}
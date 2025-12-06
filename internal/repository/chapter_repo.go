package repository

import (
	"context"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type ChapterRepo struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) *ChapterRepo {
	return &ChapterRepo{db: db}
}

func (r *ChapterRepo) Create(ctx context.Context, c *domain.Chapter) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *ChapterRepo) GetByUUID(ctx context.Context, uuid string) (*domain.Chapter, error) {
	var chapter domain.Chapter
	// Preload Slides dengan urutan sequence
	err := r.db.WithContext(ctx).
		Preload("Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Where("uuid = ?", uuid).
		First(&chapter).Error
	return &chapter, err
}

func (r *ChapterRepo) GetAllByStoryID(ctx context.Context, storyID uint) ([]domain.Chapter, error) {
	var chapters []domain.Chapter
	err := r.db.WithContext(ctx).Where("story_id = ?", storyID).Find(&chapters).Error
	return chapters, err
}

func (r *ChapterRepo) Delete(ctx context.Context, uuid string) error {
	var chapter domain.Chapter
	if err := r.db.WithContext(ctx).Select("id").Where("uuid = ?", uuid).First(&chapter).Error; err != nil {
		return err
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hapus Slide terkait
		if err := tx.Where("chapter_id = ?", chapter.ID).Delete(&domain.Slide{}).Error; err != nil {
			return err
		}
		return tx.Delete(&chapter).Error
	})
}

func (r *ChapterRepo) CreateSlide(ctx context.Context, s *domain.Slide) error {
	// Logic Transaction manual karena stored procedure 'add_slide_safe' mungkin perlu disesuaikan 
	// atau kita pakai GORM standar saja untuk update slide_count di Chapter
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(s).Error; err != nil {
			return err
		}
		// Increment Slide Count di Chapter
		return tx.Model(&domain.Chapter{}).Where("id = ?", s.ChapterID).Update("slide_count", gorm.Expr("slide_count + 1")).Error
	})
}

func (r *ChapterRepo) CountSlides(ctx context.Context, chapterID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Slide{}).Where("chapter_id = ?", chapterID).Count(&count).Error
	return count, err
}
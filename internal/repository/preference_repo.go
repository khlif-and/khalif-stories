package repository

import (
	"context"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type PreferenceRepo struct {
	db *gorm.DB
}

func NewPreferenceRepository(db *gorm.DB) *PreferenceRepo {
	return &PreferenceRepo{db: db}
}

func (r *PreferenceRepo) ClearChoices(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserChoiceStory{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserChoiceDakwah{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserChoiceHadist{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *PreferenceRepo) SaveStoryChoices(ctx context.Context, choices []domain.UserChoiceStory) error {
	if len(choices) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&choices).Error
}

func (r *PreferenceRepo) SaveDakwahChoices(ctx context.Context, choices []domain.UserChoiceDakwah) error {
	if len(choices) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&choices).Error
}

func (r *PreferenceRepo) SaveHadistChoices(ctx context.Context, choices []domain.UserChoiceHadist) error {
	if len(choices) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&choices).Error
}
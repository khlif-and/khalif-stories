package usecase

import (
	"context"

	"khalif-stories/internal/domain"

)

type PreferenceUC struct {
	repo     domain.PreferenceRepository
	catRepo  domain.CategoryRepository
}

func NewPreferenceUseCase(repo domain.PreferenceRepository, catRepo domain.CategoryRepository) *PreferenceUC {
	return &PreferenceUC{repo: repo, catRepo: catRepo}
}

func (u *PreferenceUC) SavePreferences(ctx context.Context, userID string, storyCatUUIDs, dakwahCatUUIDs, hadistCatUUIDs []string) error {
	if err := u.repo.ClearChoices(ctx, userID); err != nil {
		return err
	}

	var storyChoices []domain.UserChoiceStory
	for _, uuid := range storyCatUUIDs {
		if cat, _ := u.catRepo.GetByUUID(ctx, uuid); cat != nil {
			storyChoices = append(storyChoices, domain.UserChoiceStory{UserID: userID, CategoryID: cat.ID})
		}
	}

	var dakwahChoices []domain.UserChoiceDakwah
	for _, uuid := range dakwahCatUUIDs {
		if cat, _ := u.catRepo.GetByUUID(ctx, uuid); cat != nil {
			dakwahChoices = append(dakwahChoices, domain.UserChoiceDakwah{UserID: userID, CategoryID: cat.ID})
		}
	}

	var hadistChoices []domain.UserChoiceHadist
	for _, uuid := range hadistCatUUIDs {
		if cat, _ := u.catRepo.GetByUUID(ctx, uuid); cat != nil {
			hadistChoices = append(hadistChoices, domain.UserChoiceHadist{UserID: userID, CategoryID: cat.ID})
		}
	}

	if err := u.repo.SaveStoryChoices(ctx, storyChoices); err != nil {
		return err
	}
	if err := u.repo.SaveDakwahChoices(ctx, dakwahChoices); err != nil {
		return err
	}
	if err := u.repo.SaveHadistChoices(ctx, hadistChoices); err != nil {
		return err
	}

	return nil
}
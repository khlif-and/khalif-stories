package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	"github.com/google/uuid"

	"khalif-stories/internal/config"
	"khalif-stories/internal/domain"
	"khalif-stories/pkg/utils"

)

type ChapterUC struct {
	cfg         *config.Config
	repo        domain.ChapterRepository
	storyRepo   domain.StoryRepository
	uploader    *utils.AzureUploader
}

func NewChapterUseCase(cfg *config.Config, repo domain.ChapterRepository, storyRepo domain.StoryRepository, uploader *utils.AzureUploader) *ChapterUC {
	return &ChapterUC{cfg: cfg, repo: repo, storyRepo: storyRepo, uploader: uploader}
}

func (u *ChapterUC) Create(ctx context.Context, storyUUID string) (*domain.Chapter, error) {
	story, err := u.storyRepo.GetByUUID(ctx, storyUUID)
	if err != nil || story == nil {
		return nil, errors.New("story not found")
	}

	chapter := &domain.Chapter{
		UUID:    uuid.New().String(),
		StoryID: story.ID,
	}

	if err := u.repo.Create(ctx, chapter); err != nil {
		return nil, err
	}

	return chapter, nil
}

func (u *ChapterUC) GetByUUID(ctx context.Context, uuid string) (*domain.Chapter, error) {
	chapter, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	if chapter == nil {
		return nil, domain.ErrNotFound
	}
	return chapter, nil
}

func (u *ChapterUC) Delete(ctx context.Context, uuid string) error {
	chapter, err := u.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	for _, slide := range chapter.Slides {
		if slide.ImageURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterImages, slide.ImageURL)
		}
		if slide.SoundURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterSounds, slide.SoundURL)
		}
	}

	return u.repo.Delete(ctx, uuid)
}

func (u *ChapterUC) AddSlide(ctx context.Context, chapterUUID string, content string, sequence int, imageFile multipart.File, imageHeader *multipart.FileHeader, soundFile multipart.File, soundHeader *multipart.FileHeader) (*domain.Slide, error) {
	chapter, err := u.repo.GetByUUID(ctx, chapterUUID)
	if err != nil {
		return nil, err
	}

	count, _ := u.repo.CountSlides(ctx, chapter.ID)
	if count >= 20 {
		return nil, errors.New("maximum 20 slides per chapter reached")
	}

	var imageURL string
	if imageFile != nil {
		folderPath := ""
		url, _, err := utils.UploadAndAnalyzeImage(ctx, u.uploader, imageFile, imageHeader, u.cfg.AzureContainerChapterImages, folderPath, uuid.New().String())
		if err != nil {
			return nil, err
		}
		imageURL = url
	}

	var soundURL string
	if soundFile != nil {
		convertedFile, tempPath, err := utils.ConvertToAAC(soundFile, soundHeader.Filename)
		if err != nil {
			if imageURL != "" {
				u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterImages, imageURL)
			}
			return nil, errors.New("failed to convert audio: " + err.Error())
		}

		defer func() {
			convertedFile.Close()
			os.Remove(tempPath)
		}()

		newFilename := uuid.New().String() + ".m4a"
		folderPath := ""

		url, err := u.uploader.UploadToContainer(ctx, convertedFile, u.cfg.AzureContainerChapterSounds, folderPath+newFilename)
		if err != nil {
			if imageURL != "" {
				u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterImages, imageURL)
			}
			return nil, err
		}
		soundURL = url
	}

	slide := &domain.Slide{
		ChapterID: &chapter.ID,
		Content:   content,
		Sequence:  sequence,
		ImageURL:  imageURL,
		SoundURL:  soundURL,
	}

	if err := u.repo.CreateSlide(ctx, slide); err != nil {
		if imageURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterImages, imageURL)
		}
		if soundURL != "" {
			u.uploader.DeleteFromContainer(ctx, u.cfg.AzureContainerChapterSounds, soundURL)
		}
		return nil, err
	}

	return slide, nil
}
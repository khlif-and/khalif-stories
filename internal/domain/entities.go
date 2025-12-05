package domain

import (
	"context"
	"mime/multipart"
	"time"

)

type Category struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	UUID          string    `gorm:"type:uuid;uniqueIndex" json:"id"`
	Name          string    `gorm:"index" json:"name"`
	ImageURL      string    `json:"image_url"`
	DominantColor string    `json:"dominant_color"`
	Stories       []Story   `gorm:"foreignKey:CategoryID" json:"stories,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Story struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	UUID          string    `gorm:"type:uuid;uniqueIndex" json:"id"`
	Title         string    `gorm:"index" json:"title"`
	Description   string    `json:"description"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	DominantColor string    `json:"dominant_color"`
	CategoryID    uint      `gorm:"index" json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	UserID        string    `gorm:"index" json:"user_id"` // Field Baru: Siapa yang post
	Slides        []Slide   `gorm:"foreignKey:StoryID" json:"slides,omitempty"`
	SlideCount    int       `gorm:"default:0" json:"slide_count"`
	Status        string    `gorm:"index;default:'Draft'" json:"status"`
	CreatedAt     time.Time `gorm:"index" json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Slide struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryID   uint      `gorm:"index" json:"story_id"`
	ImageURL  string    `json:"image_url"`
	Content   string    `json:"content"`
	Sequence  int       `gorm:"index" json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByName(ctx context.Context, name string) (*Category, error)
	GetByUUID(ctx context.Context, uuid string) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Search(ctx context.Context, query string) ([]Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, uuid string) error
	UpdateColor(ctx context.Context, id uint, color string) error
}

type StoryRepository interface {
	Create(ctx context.Context, s *Story) error
	GetAll(ctx context.Context, page, limit int, sort string) ([]Story, error)
	Search(ctx context.Context, query string) ([]Story, error)
	GetByID(ctx context.Context, id uint) (*Story, error)
	GetByUUID(ctx context.Context, uuid string) (*Story, error)
	Update(ctx context.Context, s *Story) error
	UpdateColor(ctx context.Context, id uint, color string) error
	Delete(ctx context.Context, uuid string) error
	
	// TAMBAHKAN INI:
	CheckDuplicate(ctx context.Context, title, description string) (bool, error)

	CreateSlide(ctx context.Context, s *Slide) error
	CountSlides(ctx context.Context, storyID uint) (int64, error)
}

type RedisRepository interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	DeletePrefix(ctx context.Context, prefix string) error
}

type StorageRepository interface {
	Upload(file multipart.File, header *multipart.FileHeader) (string, error)
	Delete(fileURL string) error
}

type CategoryUseCase interface {
	Create(ctx context.Context, name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	Get(ctx context.Context, uuid string) (*Category, error)
	Search(ctx context.Context, query string) ([]Category, error)
	Update(ctx context.Context, uuid string, name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	Delete(ctx context.Context, uuid string) error
}

type StoryUseCase interface {
	Create(ctx context.Context, title, desc string, categoryUUID string, userID string, file multipart.File, header *multipart.FileHeader) (*Story, error)
	Update(ctx context.Context, storyUUID string, title, desc, categoryUUID, status string, file multipart.File, header *multipart.FileHeader) (*Story, error)
	GetAll(ctx context.Context, page, limit int, sort string) ([]Story, error)
	Search(ctx context.Context, query string) (*[]Story, error)
	Delete(ctx context.Context, uuid string) error
	AddSlide(ctx context.Context, storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*Slide, error)
}
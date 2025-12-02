package domain

import (
	"mime/multipart"
	"time"

)

type Category struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	UUID          string    `gorm:"type:uuid;uniqueIndex" json:"id"`
	Name          string    `json:"name"`
	ImageURL      string    `json:"image_url"`
	DominantColor string    `json:"dominant_color"`
	Stories       []Story   `gorm:"foreignKey:CategoryID" json:"stories,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Story struct {
	ID            uint      `gorm:"primaryKey" json:"-"`
	UUID          string    `gorm:"type:uuid;uniqueIndex" json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	DominantColor string    `json:"dominant_color"`
	CategoryID    uint      `json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Slides        []Slide   `gorm:"foreignKey:StoryID" json:"slides,omitempty"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Slide struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StoryID   uint      `json:"story_id"`
	ImageURL  string    `json:"image_url"`
	Content   string    `json:"content"`
	Sequence  int       `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SearchRepository interface {
	IndexCategory(category *Category) error
	SearchCategories(query string) ([]Category, error)
	DeleteCategoryIndex(uuid string) error
}

type StoryRepository interface {
	CreateCategory(category *Category) error
	GetCategoryByName(name string) (*Category, error)
	GetCategoryByUUID(uuid string) (*Category, error)
	GetAllCategories() ([]Category, error)
	UpdateCategory(category *Category) error
	DeleteCategory(uuid string) error
	UpdateCategoryColor(id uint, color string) error

	CreateStory(story *Story) error
	GetStoryByID(id uint) (*Story, error)
	GetStoryByUUID(uuid string) (*Story, error)
	GetAllStories(page, limit int, sort string) ([]Story, error)
	DeleteStory(uuid string) error
	UpdateStory(story *Story) error
	UpdateStoryColor(id uint, color string) error

	CreateSlide(slide *Slide) error
	CountSlides(storyID uint) (int64, error)
}

type StoryUseCase interface {
	CreateCategory(name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	GetAllCategories() ([]Category, error)
	GetCategory(uuid string) (*Category, error)
	UpdateCategory(uuid string, name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	DeleteCategory(uuid string) error
	SearchCategories(query string) ([]Category, error)

	CreateStory(title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*Story, error)
	GetAllStories(page, limit int, sort string) ([]Story, error)
	DeleteStory(uuid string) error

	AddSlide(storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*Slide, error)
}
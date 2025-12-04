package domain

import (
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
	Create(category *Category) error
	GetByName(name string) (*Category, error)
	GetByUUID(uuid string) (*Category, error)
	GetAll() ([]Category, error)
	Search(query string) ([]Category, error)
	Update(category *Category) error
	Delete(uuid string) error
	UpdateColor(id uint, color string) error
}

type StoryRepository interface {
	Create(story *Story) error
	GetByID(id uint) (*Story, error)
	GetByUUID(uuid string) (*Story, error)
	GetAll(page, limit int, sort string) ([]Story, error)
	Search(query string) ([]Story, error)
	Update(story *Story) error
	Delete(uuid string) error
	UpdateColor(id uint, color string) error
	CreateSlide(slide *Slide) error
	CountSlides(storyID uint) (int64, error)
}

type CategoryUseCase interface {
	Create(name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	GetAll() ([]Category, error)
	Get(uuid string) (*Category, error)
	Search(query string) ([]Category, error)
	Update(uuid string, name string, file multipart.File, header *multipart.FileHeader) (*Category, error)
	Delete(uuid string) error
}

type StoryUseCase interface {
	Create(title, desc string, categoryID uint, file multipart.File, header *multipart.FileHeader) (*Story, error)
	GetAll(page, limit int, sort string) ([]Story, error)
	Search(query string) (*[]Story, error)
	Delete(uuid string) error
	AddSlide(storyUUID string, content string, sequence int, file multipart.File, header *multipart.FileHeader) (*Slide, error)
}
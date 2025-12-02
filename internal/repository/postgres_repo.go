package repository

import (
	"errors"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

type StoryRepo struct {
	db *gorm.DB
}

func NewStoryRepository(db *gorm.DB) *StoryRepo {
	return &StoryRepo{db: db}
}

func (r *StoryRepo) GetCategoryByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Raw("SELECT * FROM categories WHERE name = ? LIMIT 1", name).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	if category.ID == 0 {
		return nil, nil
	}
	return &category, nil
}

func (r *StoryRepo) GetCategoryByUUID(uuid string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Raw("SELECT * FROM categories WHERE uuid = ? LIMIT 1", uuid).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	if category.ID == 0 {
		return nil, errors.New("category not found")
	}
	return &category, nil
}

func (r *StoryRepo) CreateCategory(c *domain.Category) error {
	query := `
		INSERT INTO categories (uuid, name, image_url, dominant_color, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW()) 
		RETURNING id, uuid, name, image_url, dominant_color, created_at, updated_at
	`
	return r.db.Raw(query, c.UUID, c.Name, c.ImageURL, c.DominantColor).Scan(c).Error
}

func (r *StoryRepo) UpdateCategory(c *domain.Category) error {
	query := `
		UPDATE categories 
		SET name = ?, image_url = ?, dominant_color = ?, updated_at = NOW() 
		WHERE id = ?
	`
	return r.db.Exec(query, c.Name, c.ImageURL, c.DominantColor, c.ID).Error
}

func (r *StoryRepo) DeleteCategory(uuid string) error {
	var category domain.Category
	// 1. Ambil ID Internal
	if err := r.db.Raw("SELECT id FROM categories WHERE uuid = ?", uuid).Scan(&category).Error; err != nil {
		return err
	}
	if category.ID == 0 {
		return errors.New("category not found")
	}

	// 2. Hard Delete Cascade Manual (Slides -> Stories -> Category)
	// Hapus semua slides dari semua stories yang ada di kategori ini
	err := r.db.Exec("DELETE FROM slides WHERE story_id IN (SELECT id FROM stories WHERE category_id = ?)", category.ID).Error
	if err != nil {
		return err
	}

	// Hapus semua stories di kategori ini
	if err := r.db.Exec("DELETE FROM stories WHERE category_id = ?", category.ID).Error; err != nil {
		return err
	}

	// 3. Hapus Category itu sendiri
	return r.db.Exec("DELETE FROM categories WHERE id = ?", category.ID).Error
}

func (r *StoryRepo) GetAllCategories() ([]domain.Category, error) {
	var cats []domain.Category
	err := r.db.Raw("SELECT * FROM categories ORDER BY id ASC").Scan(&cats).Error
	return cats, err
}

func (r *StoryRepo) UpdateCategoryColor(id uint, color string) error {
	return r.db.Exec("UPDATE categories SET dominant_color = ?, updated_at = NOW() WHERE id = ?", color, id).Error
}

func (r *StoryRepo) CreateStory(s *domain.Story) error {
	query := `
		INSERT INTO stories (uuid, title, description, thumbnail_url, dominant_color, category_id, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW()) 
		RETURNING id, uuid, title, description, thumbnail_url, dominant_color, category_id, status, created_at, updated_at
	`
	return r.db.Raw(query, s.UUID, s.Title, s.Description, s.ThumbnailURL, s.DominantColor, s.CategoryID, s.Status).Scan(s).Error
}

func (r *StoryRepo) GetAllStories(page, limit int, sort string) ([]domain.Story, error) {
	var stories []domain.Story
	offset := (page - 1) * limit

	err := r.db.Preload("Category").
		Preload("Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence ASC")
		}).
		Order(sort).
		Limit(limit).
		Offset(offset).
		Find(&stories).Error

	return stories, err
}

func (r *StoryRepo) GetStoryByID(id uint) (*domain.Story, error) {
	var story domain.Story
	err := r.db.Raw("SELECT * FROM stories WHERE id = ?", id).Scan(&story).Error
	if err != nil {
		return nil, err
	}
	return &story, nil
}

func (r *StoryRepo) GetStoryByUUID(uuid string) (*domain.Story, error) {
	var story domain.Story
	err := r.db.Raw("SELECT * FROM stories WHERE uuid = ?", uuid).Scan(&story).Error
	if err != nil {
		return nil, err
	}
	if story.ID == 0 {
		return nil, errors.New("story not found")
	}
	return &story, nil
}

func (r *StoryRepo) UpdateStory(s *domain.Story) error {
	query := `
		UPDATE stories 
		SET title = ?, description = ?, thumbnail_url = ?, dominant_color = ?, category_id = ?, status = ?, updated_at = NOW() 
		WHERE id = ?
	`
	return r.db.Exec(query, s.Title, s.Description, s.ThumbnailURL, s.DominantColor, s.CategoryID, s.Status, s.ID).Error
}

func (r *StoryRepo) UpdateStoryColor(id uint, color string) error {
	return r.db.Exec("UPDATE stories SET dominant_color = ?, updated_at = NOW() WHERE id = ?", color, id).Error
}

func (r *StoryRepo) DeleteStory(uuid string) error {
	var id uint
	if err := r.db.Raw("SELECT id FROM stories WHERE uuid = ?", uuid).Scan(&id).Error; err != nil {
		return err
	}
	if id == 0 {
		return errors.New("story not found")
	}

	if err := r.db.Exec("DELETE FROM slides WHERE story_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Exec("DELETE FROM stories WHERE id = ?", id).Error
}

func (r *StoryRepo) CreateSlide(s *domain.Slide) error {
	query := `
		INSERT INTO slides (story_id, image_url, content, sequence, created_at, updated_at) 
		VALUES (?, ?, ?, ?, NOW(), NOW()) 
		RETURNING id, story_id, image_url, content, sequence, created_at, updated_at
	`
	return r.db.Raw(query, s.StoryID, s.ImageURL, s.Content, s.Sequence).Scan(s).Error
}

func (r *StoryRepo) CountSlides(storyID uint) (int64, error) {
	var count int64
	err := r.db.Raw("SELECT count(*) FROM slides WHERE story_id = ?", storyID).Scan(&count).Error
	return count, err
}
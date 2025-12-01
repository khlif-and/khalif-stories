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

func (r *StoryRepo) GetCategoryByName(name string) (*domain.Category, error) {
	var category domain.Category
	// Cek apakah ada kategori dengan nama persis sama (Case Sensitive tergantung collation DB)
	err := r.db.Raw("SELECT * FROM categories WHERE name = ? LIMIT 1", name).Scan(&category).Error
	if err != nil {
		return nil, err
	}
	// Jika tidak ditemukan, ID akan 0
	if category.ID == 0 {
		return nil, nil
	}
	return &category, nil
}

func (r *StoryRepo) CreateCategory(c *domain.Category) error {
	query := `
		INSERT INTO categories (name, image_url, dominant_color, created_at, updated_at) 
		VALUES (?, ?, ?, NOW(), NOW()) 
		RETURNING id, name, image_url, dominant_color, created_at, updated_at
	`
	return r.db.Raw(query, c.Name, c.ImageURL, c.DominantColor).Scan(c).Error
}

func (r *StoryRepo) GetAllCategories() ([]domain.Category, error) {
	var cats []domain.Category
	err := r.db.Raw("SELECT * FROM categories ORDER BY id ASC").Scan(&cats).Error
	return cats, err
}

func (r *StoryRepo) CreateStory(s *domain.Story) error {
	query := `
		INSERT INTO stories (title, description, thumbnail_url, dominant_color, category_id, status, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW()) 
		RETURNING id, title, description, thumbnail_url, dominant_color, category_id, status, created_at, updated_at
	`
	return r.db.Raw(query, s.Title, s.Description, s.ThumbnailURL, s.DominantColor, s.CategoryID, s.Status).Scan(s).Error
}

func (r *StoryRepo) GetAllStories(page, limit int, sort string) ([]domain.Story, error) {
	var stories []domain.Story
	offset := (page - 1) * limit

	query := "SELECT * FROM stories ORDER BY " + sort + " LIMIT ? OFFSET ?"
	err := r.db.Raw(query, limit, offset).Scan(&stories).Error
	if err != nil {
		return nil, err
	}

	for i := range stories {
		r.db.Raw("SELECT * FROM categories WHERE id = ?", stories[i].CategoryID).Scan(&stories[i].Category)
		r.db.Raw("SELECT * FROM slides WHERE story_id = ? ORDER BY sequence ASC", stories[i].ID).Scan(&stories[i].Slides)
	}

	return stories, nil
}

func (r *StoryRepo) GetStoryByID(id uint) (*domain.Story, error) {
	var story domain.Story
	err := r.db.Raw("SELECT * FROM stories WHERE id = ?", id).Scan(&story).Error
	if err != nil {
		return nil, err
	}

	r.db.Raw("SELECT * FROM categories WHERE id = ?", story.CategoryID).Scan(&story.Category)
	r.db.Raw("SELECT * FROM slides WHERE story_id = ? ORDER BY sequence ASC", story.ID).Scan(&story.Slides)

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

func (r *StoryRepo) DeleteStory(id uint) error {
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
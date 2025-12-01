package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"gorm.io/gorm"

	"khalif-stories/internal/domain"

)

func SeedCategories(db *gorm.DB) {
	var count int64
	db.Raw("SELECT count(*) FROM categories").Scan(&count)

	if count > 0 {
		return
	}

	jsonFile, err := os.Open("seeds/categories.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var categories []domain.Category
	if err := json.Unmarshal(byteValue, &categories); err != nil {
		fmt.Println(err)
		return
	}

	tx := db.Begin()

	query := "INSERT INTO categories (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW()) ON CONFLICT (id) DO NOTHING"

	for _, cat := range categories {
		if err := tx.Exec(query, cat.ID, cat.Name).Error; err != nil {
			tx.Rollback()
			fmt.Println(err)
			return
		}
	}

	tx.Commit()
	fmt.Println("Seeding Categories finished")
}
package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
	_ "github.com/jackc/pgx/v5/stdlib"

)

func EnsureDBExists(dsn string) {
	dbName := extractDBName(dsn)
	rootDSN := strings.Replace(dsn, "dbname="+dbName, "dbname=postgres", 1)

	db, err := sql.Open("pgx", rootDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var exists int
	checkQuery := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s'", dbName)
	err = db.QueryRow(checkQuery).Scan(&exists)

	if err == sql.ErrNoRows {
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE \"%s\"", dbName))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ResetSchema(db *gorm.DB) {
	queries := []string{
		"DROP SCHEMA public CASCADE;",
		"CREATE SCHEMA public;",
		"GRANT ALL ON SCHEMA public TO public;",
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Fatal(err)
		}
	}
}

func extractDBName(dsn string) string {
	parts := strings.Split(dsn, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			return strings.TrimPrefix(part, "dbname=")
		}
	}
	return "khalif_stories_db"
}
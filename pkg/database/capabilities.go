package database

import (
	"log"

	"gorm.io/gorm"

)

func SetupDatabaseCapabilities(db *gorm.DB) {
	setupTriggers(db)
	setupStoredProcedures(db)
	setupRoles(db)
}

func setupTriggers(db *gorm.DB) {
	queries := []string{
		`CREATE OR REPLACE FUNCTION update_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = NOW();
			RETURN NEW;
		END;
		$$ language 'plpgsql';`,

		`DROP TRIGGER IF EXISTS update_categories_modtime ON categories;`,
		`CREATE TRIGGER update_categories_modtime BEFORE UPDATE ON categories FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();`,
		
		`DROP TRIGGER IF EXISTS update_stories_modtime ON stories;`,
		`CREATE TRIGGER update_stories_modtime BEFORE UPDATE ON stories FOR EACH ROW EXECUTE PROCEDURE update_updated_at_column();`,

		`CREATE OR REPLACE FUNCTION update_slide_count()
		RETURNS TRIGGER AS $$
		BEGIN
			IF (TG_OP = 'INSERT') THEN
				UPDATE stories SET slide_count = slide_count + 1, updated_at = NOW() WHERE id = NEW.story_id;
				RETURN NEW;
			ELSIF (TG_OP = 'DELETE') THEN
				UPDATE stories SET slide_count = slide_count - 1, updated_at = NOW() WHERE id = OLD.story_id;
				RETURN OLD;
			END IF;
			RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;`,

		`DROP TRIGGER IF EXISTS trg_update_slide_count ON slides;`,
		`CREATE TRIGGER trg_update_slide_count
		AFTER INSERT OR DELETE ON slides
		FOR EACH ROW EXECUTE PROCEDURE update_slide_count();`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Failed to execute trigger query: %v", err)
		}
	}
}

func setupStoredProcedures(db *gorm.DB) {
	queries := []string{
		`CREATE OR REPLACE PROCEDURE add_slide_safe(
			p_story_id INT, 
			p_image_url TEXT, 
			p_content TEXT, 
			p_sequence INT
		)
		LANGUAGE plpgsql
		AS $$
		DECLARE
			current_count INT;
		BEGIN
			SELECT count(*) INTO current_count FROM slides WHERE story_id = p_story_id;
			
			IF current_count >= 20 THEN
				RAISE EXCEPTION 'Maximum slide limit (20) reached for story ID %', p_story_id;
			END IF;

			INSERT INTO slides (story_id, image_url, content, sequence, created_at, updated_at)
			VALUES (p_story_id, p_image_url, p_content, p_sequence, NOW(), NOW());
		END;
		$$;`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Failed to execute stored procedure query: %v", err)
		}
	}
}

func setupRoles(db *gorm.DB) {
	queries := []string{
		`DO $$
		BEGIN
			IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'readonly_user') THEN
				CREATE ROLE readonly_user WITH LOGIN PASSWORD 'readonly_password';
			END IF;
		END
		$$;`,
		`GRANT CONNECT ON DATABASE postgres TO readonly_user;`,
		`GRANT USAGE ON SCHEMA public TO readonly_user;`,
		`GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly_user;`,
		`ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly_user;`,
	}

	for _, q := range queries {
		if err := db.Exec(q).Error; err != nil {
			log.Printf("Role setup warning (may require superuser): %v", err)
		}
	}
}
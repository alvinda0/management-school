package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"project-go/config"
)

type Migration struct {
	Version string
	Name    string
	Path    string
}

func RunMigrations() error {
	log.Println("Starting database migrations...")

	// Get all migration files
	migrations, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %v", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	// Create migrations table if not exists
	err = createMigrationsTable()
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Run each migration
	for _, migration := range migrations {
		if !isMigrationRun(migration.Version) {
			log.Printf("Running migration: %s", migration.Name)
			err = runMigration(migration)
			if err != nil {
				return fmt.Errorf("failed to run migration %s: %v", migration.Name, err)
			}
			
			err = markMigrationAsRun(migration.Version, migration.Name)
			if err != nil {
				return fmt.Errorf("failed to mark migration as run: %v", err)
			}
			
			log.Printf("Migration %s completed successfully", migration.Name)
		} else {
			log.Printf("Migration %s already run, skipping", migration.Name)
		}
	}

	log.Println("All migrations completed successfully!")
	return nil
}

func getMigrationFiles() ([]Migration, error) {
	var migrations []Migration
	
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		filename := filepath.Base(file)
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			continue
		}

		version := parts[0]
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Path:    file,
		})
	}

	return migrations, nil
}

func createMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		version VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	return config.DB.Exec(query).Error
}

func isMigrationRun(version string) bool {
	var count int64
	config.DB.Raw("SELECT COUNT(*) FROM migrations WHERE version = ?", version).Scan(&count)
	return count > 0
}

func runMigration(migration Migration) error {
	// Read migration file
	content, err := ioutil.ReadFile(migration.Path)
	if err != nil {
		return err
	}

	// Execute migration
	return config.DB.Exec(string(content)).Error
}

func markMigrationAsRun(version, name string) error {
	query := "INSERT INTO migrations (version, name) VALUES (?, ?)"
	return config.DB.Exec(query, version, name).Error
}
package postgres

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

const (
	migrationSourcePath = "file://migrations"
	migrationFilePath   = "./migrations"
)

func createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return f.Close()
}

// CreateMigrationFiles creates two new migration files with the provided filename.
// The function generates a timestamp and appends it to the filename to ensure uniqueness.
// It creates an 'up' migration file and a 'down' migration file.
// If the filename is not provided, it returns an error.
// If there is an error while creating the 'up' migration file, it returns the error.
// If there is an error while creating the 'down' migration file, it deletes the 'up' migration file and returns the error.
//
// Parameters:
// filename: The base name for the migration files. The actual file names will be prepended with a timestamp and appended with '.up.sql' or '.down.sql'.
//
// Returns:
// If successful, returns nil. If an error occurs, returns the error.
func CreateMigrationFiles(filename string) error {
	if len(filename) == 0 {
		return errors.New("migration filename is not provided")
	}

	timestamp := time.Now().Unix()
	upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.sql", migrationFilePath, timestamp, filename)
	downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.sql", migrationFilePath, timestamp, filename)

	if err := createFile(upMigrationFilePath); err != nil {
		return err
	}

	if err := createFile(downMigrationFilePath); err != nil {
		os.Remove(upMigrationFilePath)
		return err
	}

	return nil
}

// NewMigrationInstance creates a new migration instance using the provided gorm.DB data.
// It first retrieves the underlying sql.DB from the gorm.DB instance.
// If there is an error while retrieving the sql.DB, it returns the error.
// It then creates a new postgres database driver using the sql.DB and a default postgres.Config.
// If there is an error while creating the postgres driver, it returns the error.
// It then creates a new migration instance using the migration source path, the database name, and the postgres driver.
// If there is an error while creating the migration instance, it returns the error.
//
// Parameters:
// data: A pointer to a gorm.DB instance that will be used to create the migration instance.
//
// Returns:
// If successful, returns a pointer to the new migration instance and nil. If an error occurs, returns nil and the error.
func NewMigrationInstance(data *gorm.DB) (*migrate.Migrate, error) {
	db, err := data.DB()
	if err != nil {
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationSourcePath, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// RunMigration runs the provided migration.
// It calls the Up method on the migration instance, which applies all up migrations.
// If there is an error while running the migrations, it returns the error.
// If there are no migrations to apply, it returns nil.
//
// Parameters:
// migration: A pointer to a migrate.Migrate instance that will be run.
//
// Returns:
// If successful, returns nil. If an error occurs, returns the error.
func RunMigration(migration *migrate.Migrate) error {
	if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// RollbackMigration rolls back the last applied migration.
// It calls the Steps method on the migration instance with a parameter of -1, which rolls back the last migration.
// If there is an error while rolling back the migration, it returns the error.
// If there are no migrations to roll back, it returns nil.
//
// Parameters:
// migration: A pointer to a migrate.Migrate instance that will be rolled back.
//
// Returns:
// If successful, returns nil. If an error occurs, returns the error.
func RollbackMigration(migration *migrate.Migrate) error {
	if err := migration.Steps(-1); err != nil {
		return err
	}
	return nil
}

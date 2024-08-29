package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
)

type Type string

const (
	migrationSourcePath = "file://migrations"
	migrationFilePath   = "./migrations"

	PostgresDB Type = "postgres"
	MongoDB    Type = "mongodb"
)

func createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return f.Close()
}

// CreateMigrationFiles creates migration files for the specified database type.
// It generates both the up and down migration files with a timestamp and the provided filename.
// The file extensions are determined based on the database type (SQL for Postgres, JSON for MongoDB).
// If the filename is empty, it returns an error indicating that the migration filename is not found.
// If the database type is invalid, it returns an error indicating that the database type is invalid.
// If there is an error creating the up migration file, it returns the error.
// If there is an error creating the down migration file, it removes the up migration file and returns the error.
//
// Parameters:
// - filename: The name of the migration file.
// - t: The type of the database (Postgres or MongoDB).
//
// Returns:
// - error: Any error that occurred during the file creation process.
func CreateMigrationFiles(filename string, t Type) error {
	if len(filename) == 0 {
		return errors.New("migration filename is not found")
	}

	timestamp := time.Now().Unix()
	var upMigrationFilePath, downMigrationFilePath string
	switch t {
	case PostgresDB:
		upMigrationFilePath = fmt.Sprintf("%s/%d_%s.up.sql", migrationFilePath, timestamp, filename)
		downMigrationFilePath = fmt.Sprintf("%s/%d_%s.down.sql", migrationFilePath, timestamp, filename)
	case MongoDB:
		upMigrationFilePath = fmt.Sprintf("%s/%d_%s.up.json", migrationFilePath, timestamp, filename)
		downMigrationFilePath = fmt.Sprintf("%s/%d_%s.down.json", migrationFilePath, timestamp, filename)
	default:
		return errors.New("db type is invalid")
	}

	if err := createFile(upMigrationFilePath); err != nil {
		return err
	}

	if err := createFile(downMigrationFilePath); err != nil {
		os.Remove(upMigrationFilePath)
		return err
	}

	return nil
}

func newMigrationInstance(t Type, driver database.Driver) (*migrate.Migrate, error) {
	m, err := migrate.NewWithDatabaseInstance(migrationSourcePath, string(t), driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// NewPostgresInstance creates a new migration instance for Postgres.
// It initializes a Postgres driver with the provided database connection
// and then creates a new migration instance using this driver.
//
// Parameters:
// - db: A pointer to sql.DB instance representing the database connection.
//
// Returns:
// - *migrate.Migrate: A pointer to the created migration instance.
// - error: An error if the driver initialization or migration instance creation fails.
func NewPostgresInstance(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	return newMigrationInstance(PostgresDB, driver)
}

// NewMongoInstance creates a new migration instance for MongoDB.
// It initializes a MongoDB driver with the provided MongoDB client
// and then creates a new migration instance using this driver.
//
// Parameters:
// - client: A pointer to a mongo.Client instance representing the MongoDB client.
//
// Returns:
// - *migrate.Migrate: A pointer to the created migration instance.
// - error: An error if the driver initialization or migration instance creation fails.
func NewMongoInstance(client *mongo.Client) (*migrate.Migrate, error) {
	driver, err := mongodb.WithInstance(client, &mongodb.Config{TransactionMode: true})
	if err != nil {
		return nil, err
	}

	return newMigrationInstance(MongoDB, driver)
}

// RunMigration runs the provided migration.
// It calls the Up method on the migration instance, which applies all up migrations.
// If there is an error while running the migrations, it returns the error.
// If there are no migrations to apply, it returns nil.
//
// Parameters:
// migration: A pointer to migrate.Migrate instance that will be run.
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
// migration: A pointer to migrate.Migrate instance that will be rolled back.
//
// Returns:
// If successful, returns nil. If an error occurs, returns the error.
func RollbackMigration(migration *migrate.Migrate) error {
	if err := migration.Steps(-1); err != nil {
		return err
	}
	return nil
}

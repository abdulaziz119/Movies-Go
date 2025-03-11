package postgres

import (
	"Movies-Go/internal/entity"
	"Movies-Go/internal/pkg/config"
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"log"
	"os"
	"strings"
)

func NewPostgres() *bun.DB {
	dsn := "postgres://" + config.GetConf().DBUsername + ":" + config.GetConf().DBPassword + "@" +
		config.GetConf().DBHost + ":" + config.GetConf().DBPort + "/" + config.GetConf().DBName +
		"?sslmode=disable"

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := createTablesIfNotExist(db); err != nil {
		log.Printf("Error creating tables: %v", err)
	}

	if err := runMigrations(db); err != nil {
		log.Printf("Error running migrations: %v", err)
	}

	return db
}

func createTablesIfNotExist(db *bun.DB) error {
	models := []interface{}{
		(*entity.User)(nil),
		(*entity.Movie)(nil),
	}

	for _, model := range models {
		log.Printf("Creating table for model: %T", model)
		_, err := db.NewCreateTable().Model(model).IfNotExists().Exec(context.Background())
		if err != nil {
			log.Printf("Error creating table for model %T: %v", model, err)
			return err
		}
		log.Printf("Table for model %T created or already exists", model)
	}

	return nil
}

func runMigrations(db *bun.DB) error {
	migrationFiles := []string{
		"internal/pkg/script/migrations/users.sql",
		"internal/pkg/script/migrations/genres.sql",
		"internal/pkg/script/migrations/movies.sql",
	}

	for _, file := range migrationFiles {
		log.Printf("Checking migration file: %s", file)
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Migration file %s not found, skipping: %v", file, err)
			continue
		}

		log.Printf("Running migration file: %s", file)

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Error starting transaction for %s: %v", file, err)
			continue
		}

		statements := strings.Split(string(content), ";")
		executionFailed := false

		for i, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			log.Printf("Executing statement %d from %s", i+1, file)
			_, err = tx.Exec(stmt)
			if err != nil {
				log.Printf("Error executing statement %d from %s: %v", i+1, file, err)
				log.Printf("Statement was: %s", stmt)
				executionFailed = true
				break
			}
			log.Printf("Successfully executed statement %d from %s", i+1, file)
		}

		if executionFailed {
			log.Printf("Rolling back migration for %s due to errors", file)
			if err := tx.Rollback(); err != nil {
				log.Printf("Error rolling back transaction: %v", err)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Error committing transaction for %s: %v", file, err)
			continue
		}

		log.Printf("Completed migration file: %s", file)
	}

	log.Println("All migrations completed")
	return nil
}

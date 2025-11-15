package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/pasarly/backend/services/auth-service/configs"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(cfg *configs.Database, path string) (*Migrator, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database for migration: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, cfg.Name, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
	if err := m.migrate.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	return nil
}

func (m *Migrator) Down(steps int) error {
	var err error
	if steps == 0 {
		err = m.migrate.Down()
	} else {
		err = m.migrate.Steps(-steps)
	}

	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}

func (m *Migrator) Close() error {
	errS, errDB := m.migrate.Close()
	if errS != nil {
		return fmt.Errorf("failed to close migration source: %w", errS)
	}
	if errDB != nil {
		return fmt.Errorf("failed to close migration database: %w", errDB)
	}
	return nil
}

package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ritchieridanko/pasarly/backend/services/user/configs"
)

type Migrator struct {
	config  *configs.Database
	migrate *migrate.Migrate
}

func NewMigrator(cfg *configs.Database, path string) (*Migrator, error) {
	db, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %w", err)
	}

	migrate, err := migrate.NewWithDatabaseInstance("file://"+path, cfg.Name, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migrator: %w", err)
	}

	return &Migrator{config: cfg, migrate: migrate}, nil
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
	es, ed := m.migrate.Close()
	if es != nil {
		return fmt.Errorf("failed to close migration source: %w", es)
	}
	if ed != nil {
		return fmt.Errorf("failed to close migration database: %w", ed)
	}
	return nil
}

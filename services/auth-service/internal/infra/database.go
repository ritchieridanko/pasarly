package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/pasarly/auth-service/configs"
)

type Database struct {
	database *pgxpool.Pool
}

func NewDatabase(cfg *configs.Database, l *Logger) (*Database, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to build pgxpool config: %w", err)
	}

	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	l.Sugar().Infof("✅ [DATABASE] initialized (host=%s, port=%d, name=%s)", cfg.Host, cfg.Port, cfg.Name)
	return &Database{database: db}, nil
}

func (d *Database) DB() *pgxpool.Pool {
	return d.database
}

func (d *Database) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := d.database.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (d *Database) Close() {
	d.database.Close()
}

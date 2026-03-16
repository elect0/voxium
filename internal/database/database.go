package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"

	"github.com/elect0/voxium/internal/config"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func Init(ctx context.Context, cfg config.DatabaseConfig, log *zap.Logger) (*pgxpool.Pool, error) {

	migrationDB, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database for migration :%w", err)
	}

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationFS,
		Root:       "migrations",
	}

	n, err := migrate.Exec(migrationDB, "postgres", migrations, migrate.Up)
	if err != nil {
		migrationDB.Close()

		return nil, fmt.Errorf("Migration failed: %w", err)
	}

	if n > 0 {
		log.Info("Applied Migrations", zap.Int("Count", n))
	} else {
		log.Info("Schema is up to date!")
	}

	migrationDB.Close()

	config, err := pgxpool.ParseConfig(cfg.URL)

	if err != nil {
		return nil, fmt.Errorf("Unable to parse database config: %w", err)
	}

	// Remove hard-coded values
	config.MaxConns = 25
	config.MinConns = 2
	config.MaxConnLifetime = 1 * time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("Unable to create connection pool: %w", err)
	}

	return pool, nil
}

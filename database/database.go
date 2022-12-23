package database

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Open(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Ping(ctx context.Context, db *sql.DB) error {
	var pingErr error

	for i := 1; ; i++ {
		pingErr = db.PingContext(ctx)
		if pingErr == nil {
			break
		}

		time.Sleep(time.Duration(i) * 100 * time.Millisecond)

		if ctx.Err() != nil {
			return ctx.Err()
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	const query = `SELECT true`
	var res bool
	return db.QueryRowContext(ctx, query).Scan(&res)
}

func Migrate(ctx context.Context, db *sql.DB) error {
	if err := Ping(ctx, db); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("source driver: %w", err)
	}

	databaseDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("database driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgre", databaseDriver)
	if err != nil {
		return fmt.Errorf("migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate :%w", err)
		}
	}

	return nil
}

package db

import (
	"context"
	"fmt"
	"time"

	"github.com/banggibima/be-assignment/config"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func Init(cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := fmt.Sprint("postgres://", cfg.Postgres.User, ":", cfg.Postgres.Password, "@", cfg.Postgres.Host, ":", cfg.Postgres.Port, "/", cfg.Postgres.Database, "?sslmode=", cfg.Postgres.SSLMode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

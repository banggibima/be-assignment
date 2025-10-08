package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseTransactionRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseTransactionRepository(
	db *pgxpool.Pool,
) *DatabaseTransactionRepository {
	return &DatabaseTransactionRepository{
		db: db,
	}
}

func (r *DatabaseTransactionRepository) FetchBatch(ctx context.Context, from, to string, limit, offset int) (pgx.Rows, error) {
	query := "SELECT id, order_id, merchant_id, amount, fee, status, paid_at, created_at, updated_at "
	query += "FROM transactions WHERE paid_at BETWEEN $1 AND $2 "
	query += "ORDER BY paid_at ASC LIMIT $3 OFFSET $4"

	return r.db.Query(ctx, query, from, to, limit, offset)
}

func (r *DatabaseTransactionRepository) CountByDateRange(ctx context.Context, from, to string) (int, error) {
	query := "SELECT COUNT(*) FROM transactions WHERE paid_at BETWEEN $1 AND $2"
	row := r.db.QueryRow(ctx, query, from, to)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

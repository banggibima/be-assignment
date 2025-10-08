package repositories

import (
	"context"
	"errors"

	"github.com/banggibima/be-assignment/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type DatabaseProductRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseProductRepository(db *pgxpool.Pool) *DatabaseProductRepository {
	return &DatabaseProductRepository{db: db}
}

func (r *DatabaseProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := "SELECT id, name, stock, price, created_at, updated_at "
	query += "FROM products WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)

	var p models.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Stock, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *DatabaseProductRepository) UpdateStock(ctx context.Context, tx pgx.Tx, id string, qty int) (bool, error) {
	query := "UPDATE products "
	query += "SET stock = stock - $1, updated_at = NOW() "
	query += "WHERE id = $2 AND stock >= $1 "
	query += "RETURNING stock"

	var newStock int
	err := tx.QueryRow(ctx, query, qty, id).Scan(&newStock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

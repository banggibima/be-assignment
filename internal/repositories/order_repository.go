package repositories

import (
	"context"

	"github.com/banggibima/be-assignment/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseOrderRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseOrderRepository(db *pgxpool.Pool) *DatabaseOrderRepository {
	return &DatabaseOrderRepository{db: db}
}

func (r *DatabaseOrderRepository) Create(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	order.ID = uuid.New().String()

	query := "INSERT INTO orders (id, product_id, buyer_id, quantity, total_price, created_at, updated_at) "
	query += "VALUES ($1, $2, $3, $4, $5, NOW(), NOW())"

	_, err := tx.Exec(ctx, query, order.ID, order.ProductID, order.BuyerID, order.Quantity, order.TotalPrice)
	return err
}

func (r *DatabaseOrderRepository) GetByID(ctx context.Context, id string) (*models.Order, error) {
	query := "SELECT id, product_id, buyer_id, quantity, total_price, created_at, updated_at "
	query += "FROM orders WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)

	var o models.Order
	if err := row.Scan(&o.ID, &o.ProductID, &o.BuyerID, &o.Quantity, &o.TotalPrice, &o.CreatedAt, &o.UpdatedAt); err != nil {
		return nil, err
	}
	return &o, nil
}

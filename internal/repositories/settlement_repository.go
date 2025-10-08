package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseSettlementRepository struct {
	db *pgxpool.Pool
}

func NewDatabaseSettlementRepository(
	db *pgxpool.Pool,
) *DatabaseSettlementRepository {
	return &DatabaseSettlementRepository{
		db: db,
	}
}

func (r *DatabaseSettlementRepository) UpsertJob(ctx context.Context, tx pgx.Tx, runID, merchantID, date string, gross, fee, net, txnCount int) error {
	id := uuid.New().String()

	query := "INSERT INTO settlements (id, merchant_id, date, gross_amount, fee_amount, net_amount, txn_count, unique_run_id, created_at, updated_at) "
	query += "VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) "
	query += "ON CONFLICT (merchant_id, date) "
	query += "DO UPDATE SET "
	query += "gross_amount = settlements.gross_amount + EXCLUDED.gross_amount, "
	query += "fee_amount = settlements.fee_amount + EXCLUDED.fee_amount, "
	query += "net_amount = settlements.net_amount + EXCLUDED.net_amount, "
	query += "txn_count = settlements.txn_count + EXCLUDED.txn_count, "
	query += "unique_run_id = EXCLUDED.unique_run_id, "
	query += "updated_at = NOW()"

	_, err := tx.Exec(ctx, query, id, merchantID, date, gross, fee, net, txnCount, runID)
	return err
}

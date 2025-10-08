package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	FetchBatch(ctx context.Context, from, to string, limit, offset int) (pgx.Rows, error)
	CountByDateRange(ctx context.Context, from, to string) (int, error)
}

type SettlementService struct {
	db              *pgxpool.Pool
	transactionRepo TransactionRepository
}

func NewSettlementService(
	db *pgxpool.Pool,
	transactionRepo TransactionRepository,
) *SettlementService {
	return &SettlementService{
		db:              db,
		transactionRepo: transactionRepo,
	}
}

func (s *SettlementService) ProcessSettlementWithCancel(
	ctx context.Context,
	from, to string,
	cancel <-chan struct{},
	onProgress func(processed, total int),
) error {
	total, err := s.transactionRepo.CountByDateRange(ctx, from, to)
	if err != nil {
		return err
	}

	offset := 0
	limit := 5000
	processed := 0

	for {
		select {
		case <-cancel:
			fmt.Println("[Settlement] Job cancelled gracefully.")
			return context.Canceled
		default:
		}

		rows, err := s.transactionRepo.FetchBatch(ctx, from, to, limit, offset)
		if err != nil {
			return err
		}

		count := 0
		for rows.Next() {
			count++
			processed++

			if onProgress != nil {
				onProgress(processed, total)
			}
		}
		rows.Close()

		if count < limit {
			break
		}
		offset += limit
	}

	return nil
}

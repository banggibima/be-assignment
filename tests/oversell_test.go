package tests

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/banggibima/be-assignment/internal/dto"
	"github.com/banggibima/be-assignment/internal/repositories"
	"github.com/banggibima/be-assignment/internal/services"
	"github.com/jackc/pgx/v5/pgxpool"
)

func mustConnectDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@12345678:5432/be_assignment?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("cannot parse config: %v", err)
	}

	// pool cukup untuk 500 goroutine
	config.MaxConns = 50

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("cannot connect to database: %v", err)
	}
	return pool
}

func TestNoOversell(t *testing.T) {
	pool := mustConnectDB(t)
	defer pool.Close()
	ctx := context.Background()

	// Reset product stock
	_, err := pool.Exec(ctx, `
		INSERT INTO products (id, name, stock, price, created_at, updated_at)
		VALUES ('product-1', 'Limited Product', 100, 1000, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET stock = 100, price = 1000, updated_at = NOW();
	`)
	if err != nil {
		t.Fatalf("failed to seed product: %v", err)
	}

	// Bersihkan orders
	_, _ = pool.Exec(ctx, `DELETE FROM orders WHERE product_id = 'product-1'`)

	productRepo := repositories.NewDatabaseProductRepository(pool)
	orderRepo := repositories.NewDatabaseOrderRepository(pool)
	orderService := services.NewOrderService(pool, orderRepo, productRepo)

	const totalBuyers = 500
	var successCount int32
	var failCount int32

	// Precompute unique BuyerIDs
	buyerIDs := make([]string, totalBuyers)
	for i := 0; i < totalBuyers; i++ {
		buyerIDs[i] = fmt.Sprintf("buyer-%03d", i)
	}

	var wg sync.WaitGroup
	wg.Add(totalBuyers)

	start := time.Now()

	for i := 0; i < totalBuyers; i++ {
		go func(i int) {
			defer wg.Done()
			req := dto.CreateOrderRequest{
				ProductID: "product-1",
				BuyerID:   buyerIDs[i],
				Quantity:  1,
			}

			_, err := orderService.CreateOrder(ctx, req)
			if err != nil {
				if err == services.ErrOutOfStock {
					atomic.AddInt32(&failCount, 1)
				} else {
					t.Errorf("unexpected error for buyer %s: %v", buyerIDs[i], err)
					atomic.AddInt32(&failCount, 1)
				}
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	t.Logf("Test finished: success=%d, fail=%d, duration=%s", successCount, failCount, duration)

	// Pastikan tidak oversell
	if successCount > 100 {
		t.Fatalf("oversell detected! success=%d (expected max 100)", successCount)
	}

	// Check remaining stock
	var stock int
	err = pool.QueryRow(ctx, `SELECT stock FROM products WHERE id='product-1'`).Scan(&stock)
	if err != nil {
		t.Fatalf("failed to get stock: %v", err)
	}

	expectedStock := 100 - int(successCount)
	if stock != expectedStock {
		t.Fatalf("inconsistent stock! got %d, expected %d", stock, expectedStock)
	}

	t.Logf("Final stock: %d (OK)", stock)
}

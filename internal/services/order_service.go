package services

import (
	"context"
	"errors"

	"github.com/banggibima/be-assignment/internal/dto"
	"github.com/banggibima/be-assignment/internal/models"
	"github.com/banggibima/be-assignment/internal/repositories"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOutOfStock      = errors.New("OUT_OF_STOCK")
	ErrProductNotFound = errors.New("PRODUCT_NOT_FOUND")
)

type OrderRepository interface {
	Create(ctx context.Context, tx pgx.Tx, order *models.Order) error
	GetByID(ctx context.Context, id string) (*models.Order, error)
}

type ProductRepository interface {
	GetByID(ctx context.Context, id string) (*models.Product, error)
	UpdateStock(ctx context.Context, tx pgx.Tx, id string, qty int) (bool, error)
}

type OrderService struct {
	db          *pgxpool.Pool
	orderRepo   OrderRepository
	productRepo ProductRepository
}

func NewOrderService(
	db *pgxpool.Pool,
	orderRepo OrderRepository,
	productRepo ProductRepository,
) *OrderService {
	return &OrderService{
		db:          db,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.CreateOrderResponse, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	updated, err := s.productRepo.UpdateStock(ctx, tx, req.ProductID, req.Quantity)
	if err != nil {
		return nil, err
	}
	if !updated {
		return nil, ErrOutOfStock
	}

	orderID := uuid.New().String()
	total := 1000 * req.Quantity
	err = s.orderRepo.Create(ctx, tx, &models.Order{
		ID:         orderID,
		ProductID:  req.ProductID,
		BuyerID:    req.BuyerID,
		Quantity:   req.Quantity,
		TotalPrice: total,
	})
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	res := &dto.CreateOrderResponse{
		ID:         orderID,
		ProductID:  req.ProductID,
		BuyerID:    req.BuyerID,
		Quantity:   req.Quantity,
		TotalPrice: total,
		Status:     "SUCCESS",
	}

	return res, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, id string) (*dto.GetOrderResponse, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrNotFound) {
			return nil, repositories.ErrNotFound
		}
		return nil, err
	}

	res := &dto.GetOrderResponse{
		ID:         order.ID,
		ProductID:  order.ProductID,
		BuyerID:    order.BuyerID,
		Quantity:   order.Quantity,
		TotalPrice: order.TotalPrice,
	}

	return res, nil
}

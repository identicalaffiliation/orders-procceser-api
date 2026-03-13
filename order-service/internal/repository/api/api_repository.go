package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/storage"
	"github.com/jmoiron/sqlx"
)

var (
	ErrEmptyItems error = errors.New("empty items slice")
	ErrEmptyOrder error = errors.New("empty order")
)

type ItemsRepository interface {
	CreateItems(ctx context.Context, tx *sqlx.Tx, items []*models.Item) ([]*models.Item, error)
}

type OrdersRepository interface {
	CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) (*models.Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
}

type apiRepository struct {
	db               *storage.Postgres
	itemsRepository  ItemsRepository
	ordersRepository OrdersRepository
	logger           logger.Logger
}

func NewAPIRepository(d *storage.Postgres, ir ItemsRepository, or OrdersRepository, l logger.Logger) *apiRepository {
	return &apiRepository{db: d, itemsRepository: ir, ordersRepository: or, logger: l}
}

func (r *apiRepository) CreateOrderAndItems(ctx context.Context, orderWithItems *models.OrderWithItems) (*models.Order, error) {
	tx, err := r.db.DB.BeginTxx(ctx, nil)
	if err != nil {
		r.logger.Error("begin tx", "error", err)
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			r.logger.Error("tx rollback", "error", err)
		}
	}()

	items, err := r.itemsRepository.CreateItems(ctx, tx, orderWithItems.Items)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		r.logger.Error("empty items slice", "error", ErrEmptyItems)
		return nil, ErrEmptyItems
	}

	order, err := r.ordersRepository.CreateOrder(ctx, tx, orderWithItems.Order)
	if err != nil {
		return nil, err
	}

	if order == nil {
		r.logger.Error("empty order struct", "error", err)
		return nil, ErrEmptyOrder
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("commit tx", "error", err)
		return nil, fmt.Errorf("commit tx: %w", err)
	}

	r.logger.Debug("order created in database")

	return order, nil
}

func (r *apiRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	return r.ordersRepository.GetOrder(ctx, orderID)
}

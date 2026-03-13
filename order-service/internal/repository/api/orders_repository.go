package api

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/storage"
	"github.com/jmoiron/sqlx"
)

type ordersRepository struct {
	db     *storage.Postgres
	logger logger.Logger
}

func NewOrdersRepository(db *storage.Postgres, l logger.Logger) *ordersRepository {
	return &ordersRepository{db: db, logger: l}
}

func (r *ordersRepository) CreateOrder(ctx context.Context, tx *sqlx.Tx, order *models.Order) (*models.Order, error) {
	if err := tx.QueryRowxContext(ctx, CREATE_ORDER, order.ID, order.Status, order.TotalPrice,
		order.TotalQuantity).Scan(&order.ID, &order.Created); err != nil {
		r.logger.Error("create order", "error", err)
		return nil, fmt.Errorf("create order: %w", err)
	}

	return order, nil
}

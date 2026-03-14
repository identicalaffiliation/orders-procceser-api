package api

import (
	"context"
	"fmt"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/storage"
	"github.com/jmoiron/sqlx"
)

type itemsRepository struct {
	db     *storage.Postgres
	logger logger.Logger
}

func newItemsRepository(db *storage.Postgres, l logger.Logger) *itemsRepository {
	return &itemsRepository{db: db, logger: l}
}

func (r *itemsRepository) createItem(ctx context.Context, tx *sqlx.Tx, item *models.Item) (*models.Item, error) {
	if err := tx.QueryRowxContext(ctx, CREATE_ITEM, item.ID, item.OrderID, item.Title,
		item.Price, item.Quantity).Scan(&item.ID, &item.Created); err != nil {
		r.logger.Error("create item", "error", err)
		return nil, fmt.Errorf("create item: %w", err)
	}

	return item, nil
}

func (r *itemsRepository) CreateItems(ctx context.Context, tx *sqlx.Tx, items []*models.Item) ([]*models.Item, error) {
	returned := make([]*models.Item, 0, len(items))

	for _, item := range items {
		val, err := r.createItem(ctx, tx, item)
		if err != nil {
			return nil, err
		}

		returned = append(returned, val)
	}

	return returned, nil
}

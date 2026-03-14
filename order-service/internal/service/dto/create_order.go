package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
)

type CreateOrderRequest struct {
	Items []items `json:"items" validate:"required"`
}

type items struct {
	Title    string  `json:"title" validate:"required"`
	Price    float64 `json:"price" validate:"required,gt=0"`
	Quantity int     `json:"quantity" validate:"required,min=1"`
}

type OrderResponse struct {
	ID      uuid.UUID     `json:"id"`
	Status  models.Status `json:"status"`
	Created time.Time     `json:"created"`
}

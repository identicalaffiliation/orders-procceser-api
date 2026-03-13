package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
)

type CreateOrderRequest struct {
	Items []items `json:"items"`
}

type items struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderResponse struct {
	ID      uuid.UUID     `json:"id"`
	Status  models.Status `json:"status"`
	Created time.Time     `json:"created"`
}

package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/repository/models"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/service/dto"
)

type APIRepository interface {
	CreateOrderAndItems(ctx context.Context, orderWithItems *models.OrderWithItems) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
}

type orderService struct {
	repo   APIRepository
	logger logger.Logger
}

func NewOrderService(r APIRepository, l logger.Logger) *orderService {
	return &orderService{repo: r, logger: l}
}

func (s *orderService) CreateOrder(ctx context.Context, request *dto.CreateOrderRequest,
) (*dto.OrderResponse, error) {
	var (
		totalPrice    float64
		totalQuantity int
	)

	orderID := uuid.New()
	items := make([]*models.Item, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, &models.Item{
			ID:       uuid.New(),
			OrderID:  orderID,
			Title:    item.Title,
			Price:    item.Price,
			Quantity: item.Quantity,
		})

		totalPrice += item.Price * float64(item.Quantity)
		totalQuantity += item.Quantity
	}

	order, err := s.repo.CreateOrderAndItems(ctx, &models.OrderWithItems{
		Order: &models.Order{
			ID:            orderID,
			Status:        models.STATUS_CREATED,
			TotalPrice:    totalPrice,
			TotalQuantity: totalQuantity,
		},
		Items: items,
	})
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		ID:      order.ID,
		Status:  order.Status,
		Created: order.Created,
	}, nil
}

func (s *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*dto.OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &dto.OrderResponse{
		ID:      order.ID,
		Status:  order.Status,
		Created: order.Created,
	}, nil
}

package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/config"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/logger"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/service/dto"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/pkg/validation"
	"github.com/labstack/echo"
)

type OrderService interface {
	CreateOrder(ctx context.Context, request *dto.CreateOrderRequest) (*dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*dto.OrderResponse, error)
}

type OrdersAPI struct {
	service   OrderService
	logger    logger.Logger
	validator validation.Validator
	config    *config.ServiceConfig
	engine    *echo.Echo
}

func NewOrdersAPI(s OrderService, l logger.Logger, c *config.ServiceConfig) *OrdersAPI {
	return &OrdersAPI{
		service:   s,
		logger:    l,
		validator: validation.NewValidator(),
		config:    c,
	}
}

func (api *OrdersAPI) SetupAPI() {
	api.engine = echo.New()

	address := fmt.Sprintf("%s:%s", api.config.ServerConfig.Host, api.config.ServerConfig.Port)
	api.engine.Server.Addr = address

	api.engine.POST("/api/v1/orders", api.createOrder)
}

func (api *OrdersAPI) StartAPI() error {
	return api.engine.Start(api.engine.Server.Addr)
}

func (api *OrdersAPI) ShutdownAPI(ctx context.Context) error {
	return api.engine.Shutdown(ctx)
}

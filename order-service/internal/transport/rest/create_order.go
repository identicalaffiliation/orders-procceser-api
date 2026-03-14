package rest

import (
	"net/http"

	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/service/dto"
	"github.com/labstack/echo"
)

func (api *OrdersAPI) createOrder(ctx echo.Context) error {
	var request dto.CreateOrderRequest
	if err := ctx.Bind(&request); err != nil {
		return echo.ErrBadRequest
	}

	if err := api.validator.CreateOrderRequestValidate(&request); err != nil {
		return echo.ErrBadRequest
	}

	response, err := api.service.CreateOrder(ctx.Request().Context(), &request)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, response)
}

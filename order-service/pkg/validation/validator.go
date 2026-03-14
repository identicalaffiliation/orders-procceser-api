package validation

import (
	"fmt"

	validate "github.com/go-playground/validator"
	"github.com/identicalaffiliation/orders-procceser-api/order-service/internal/service/dto"
)

type Validator interface {
	CreateOrderRequestValidate(req *dto.CreateOrderRequest) error
}

type validator struct {
	v *validate.Validate
}

func NewValidator() Validator {
	return &validator{v: validate.New()}
}

func (v *validator) CreateOrderRequestValidate(req *dto.CreateOrderRequest) error {
	if err := v.v.Struct(req); err != nil {
		return fmt.Errorf("validate order: %w", err)
	}

	for _, item := range req.Items {
		if err := v.v.Struct(item); err != nil {
			return fmt.Errorf("validate item: %w", err)
		}
	}

	return nil
}

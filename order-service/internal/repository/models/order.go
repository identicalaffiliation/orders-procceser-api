package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	STATUS_CREATED    Status = "created"
	STATUS_RESOLVED   Status = "resolved"
	STATUS_PAID       Status = "paid"
	STATUS_DELIVERING Status = "delivering"
)

type Status string

type Order struct {
	ID            uuid.UUID `db:"id"`
	Status        Status    `db:"status"`
	TotalPrice    float64   `db:"total_price"`
	TotalQuantity int       `db:"total_quantity"`
	Created       time.Time `db:"created"`
	Updated       time.Time `db:"updated"`
}

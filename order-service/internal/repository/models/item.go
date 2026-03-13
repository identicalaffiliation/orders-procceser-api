package models

import (
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID       uuid.UUID `db:"id"`
	OrderID  uuid.UUID `db:"order_id"`
	Title    string    `db:"title"`
	Price    float64   `db:"price"`
	Quantity int       `db:"quantity"`
	Created  time.Time `db:"created"`
	Updated  time.Time `db:"updated"`
}

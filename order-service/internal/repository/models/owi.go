package models

type OrderWithItems struct {
	Order *Order
	Items []*Item
}

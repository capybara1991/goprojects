package entity

import "errors"

var (
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrProductNotFound       = errors.New("product not found")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrInvalidOrderStatus    = errors.New("invalid order status")
	ErrEmptyOrderItems       = errors.New("order items is empty")
)

type OrderItem struct {
	SKU   uint32
	Count uint32
}

type Status int

const (
	StatusUnknown Status = iota
	StatusAwaitingPayment
	StatusPaid
	StatusCancelled
)

type Order struct {
	ID     uint32
	UserID int64
	Items  []OrderItem
	Status Status
}

type Product struct {
	Name  string
	Price uint64
}

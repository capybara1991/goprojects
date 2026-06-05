package entity

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrCartNotFound      = errors.New("cart not found")
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrCartIsEmpty       = errors.New("cart is empty")
)

type CartItem struct {
	SKU   uint32
	Count uint32
}

type ProductInfo struct {
	Name  string
	Price uint64
}

type ListCartItem struct {
	SKU   uint32
	Count uint32
	Name  string
	Price uint64
}

package controller

import (
	"fmt"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/loms/entity"
)

func validateCreateOrderRequest(req *inventorypb.CreateOrderRequest) error {
	if req.UserId <= 0 {
		return fmt.Errorf("invalid user id")
	}
	if len(req.Item) == 0 {
		return entity.ErrEmptyOrderItems
	}
	for _, item := range req.Item {
		if item.Sku == 0 {
			return fmt.Errorf("invalid item sku")
		}
		if item.Count == 0 {
			return fmt.Errorf("invalid item count")
		}
	}
	return nil
}

func validateOrderID(orderID int64) error {
	if orderID <= 0 {
		return entity.ErrOrderNotFound
	}
	return nil
}

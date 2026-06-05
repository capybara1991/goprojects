package controller

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/loms/entity"
)

func (a *API) CreateOrder(ctx context.Context, req *inventorypb.CreateOrderRequest) (*inventorypb.CreateOrderResponse, error) {
	if err := validateCreateOrderRequest(req); err != nil {
		return nil, mapLOMSError(err)
	}

	items := make([]entity.OrderItem, 0, len(req.Item))
	for _, item := range req.Item {
		items = append(items, entity.OrderItem{SKU: item.Sku, Count: item.Count})
	}

	orderID, err := a.orderService.CreateOrder(ctx, req.UserId, items)
	if err != nil {
		return nil, mapLOMSError(err)
	}

	return &inventorypb.CreateOrderResponse{OrderId: int64(orderID)}, nil
}

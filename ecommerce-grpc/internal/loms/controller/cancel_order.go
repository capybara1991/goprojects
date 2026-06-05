package controller

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
)

func (a *API) CancelOrder(ctx context.Context, req *inventorypb.CancelOrderRequest) (*inventorypb.CancelOrderResponse, error) {
	if err := validateOrderID(req.OrderId); err != nil {
		return nil, mapLOMSError(err)
	}

	if err := a.orderService.CancelOrder(ctx, uint32(req.OrderId)); err != nil {
		return nil, mapLOMSError(err)
	}

	return &inventorypb.CancelOrderResponse{Status: inventorypb.OrderStatus_ORDER_STATUS_CANCELLED}, nil
}

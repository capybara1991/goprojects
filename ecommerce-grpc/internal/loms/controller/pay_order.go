package controller

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
)

func (a *API) PayOrder(ctx context.Context, req *inventorypb.PayOrderRequest) (*inventorypb.PayOrderResponse, error) {
	if err := validateOrderID(req.OrderId); err != nil {
		return nil, mapLOMSError(err)
	}

	if err := a.orderService.PayOrder(ctx, uint32(req.OrderId)); err != nil {
		return nil, mapLOMSError(err)
	}

	return &inventorypb.PayOrderResponse{Status: inventorypb.OrderStatus_ORDER_STATUS_PAID}, nil
}

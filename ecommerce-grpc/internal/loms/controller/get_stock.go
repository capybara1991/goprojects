package controller

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetStock(ctx context.Context, req *inventorypb.GetStockRequest) (*inventorypb.GetStockResponse, error) {
	if req.Sku == 0 {
		return nil, status.Error(codes.InvalidArgument, "sku must be positive")
	}

	count, err := a.stocksRepo.GetStocks(ctx, req.Sku)
	if err != nil {
		return nil, mapLOMSError(err)
	}

	return &inventorypb.GetStockResponse{Count: uint32(count)}, nil
}

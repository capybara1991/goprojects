package controller

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/loms/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetProductInfo(ctx context.Context, req *inventorypb.GetProductInfoRequest) (*inventorypb.GetProductInfoResponse, error) {
	if req.Sku == 0 {
		return nil, status.Error(codes.InvalidArgument, "sku must be positive")
	}

	product, ok := a.products[req.Sku]
	if !ok {
		return nil, mapLOMSError(entity.ErrProductNotFound)
	}

	return &inventorypb.GetProductInfoResponse{Name: product.Name, Price: product.Price}, nil
}

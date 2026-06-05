package controller

import (
	"context"
	"errors"

	cartpb "grpc_pr/gen/cart"
	"grpc_pr/internal/cart/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) AddItem(ctx context.Context, req *cartpb.AddItemRequest) (*cartpb.AddItemResponse, error) {
	if err := validateAddItemRequest(req); err != nil {
		return nil, err
	}

	if err := a.itemService.AddItem(ctx, req.UserId, req.Sku, req.Count); err != nil {
		switch {
		case errors.Is(err, entity.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, "product not found")
		case errors.Is(err, entity.ErrInsufficientStock):
			return nil, status.Error(codes.FailedPrecondition, "insufficient stock")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cartpb.AddItemResponse{Message: "item added to cart"}, nil
}

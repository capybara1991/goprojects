package controller

import (
	"context"
	"errors"

	cartpb "grpc_pr/gen/cart"
	"grpc_pr/internal/cart/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) DeleteItem(ctx context.Context, req *cartpb.DeleteItemRequest) (*cartpb.DeleteItemResponse, error) {
	if err := validateDeleteItemRequest(req); err != nil {
		return nil, err
	}

	if err := a.itemService.DeleteItem(ctx, req.UserId, req.Sku); err != nil {
		resp := &cartpb.DeleteItemResponse{Status: cartpb.DeleteItemStatus_NOTDELETED}
		switch {
		case errors.Is(err, entity.ErrCartNotFound), errors.Is(err, entity.ErrCartItemNotFound):
			return resp, status.Error(codes.NotFound, err.Error())
		default:
			return resp, status.Error(codes.Internal, err.Error())
		}
	}

	return &cartpb.DeleteItemResponse{Status: cartpb.DeleteItemStatus_DELETED}, nil
}

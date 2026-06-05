package controller

import (
	"context"
	"errors"

	cartpb "grpc_pr/gen/cart"
	"grpc_pr/internal/cart/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) Checkout(ctx context.Context, req *cartpb.CheckoutRequest) (*cartpb.CheckoutResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}

	orderID, err := a.itemService.Checkout(ctx, req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrCartNotFound), errors.Is(err, entity.ErrCartIsEmpty):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, entity.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &cartpb.CheckoutResponse{OrderId: orderID}, nil
}

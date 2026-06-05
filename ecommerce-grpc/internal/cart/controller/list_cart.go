package controller

import (
	"context"
	"errors"

	cartpb "grpc_pr/gen/cart"
	"grpc_pr/internal/cart/entity"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) ListCart(ctx context.Context, req *cartpb.ListCartRequest) (*cartpb.ListCartResponse, error) {
	if err := validateUserID(req.UserId); err != nil {
		return nil, err
	}

	items, totalPrice, err := a.itemService.ListCart(ctx, req.UserId)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrCartIsEmpty):
			return &cartpb.ListCartResponse{Items: nil, TotalPrice: 0}, nil
		case errors.Is(err, entity.ErrCartNotFound), errors.Is(err, entity.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	respItems := make([]*cartpb.CartItem, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, &cartpb.CartItem{
			Sku:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return &cartpb.ListCartResponse{Items: respItems, TotalPrice: totalPrice}, nil
}

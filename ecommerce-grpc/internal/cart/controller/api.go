package controller

import (
	cartpb "grpc_pr/gen/cart"
	"grpc_pr/internal/cart/service"
)

type API struct {
	cartpb.UnimplementedCartServiceServer

	itemService *service.ItemService
}

func NewAPI(itemService *service.ItemService) *API {
	return &API{itemService: itemService}
}

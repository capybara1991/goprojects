package controller

import (
	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/loms/entity"
	"grpc_pr/internal/loms/repository"
	"grpc_pr/internal/loms/service"
)

type API struct {
	inventorypb.UnimplementedInventoryServiceServer

	orderService *service.OrderService
	stocksRepo   repository.StocksRepository
	products     map[uint32]entity.Product
}

func NewAPI(orderService *service.OrderService, stocksRepo repository.StocksRepository, products map[uint32]entity.Product) *API {
	return &API{orderService: orderService, stocksRepo: stocksRepo, products: products}
}

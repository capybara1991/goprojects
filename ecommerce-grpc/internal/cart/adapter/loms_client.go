package adapter

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/cart/entity"
)

type LOMSClient struct {
	client inventorypb.InventoryServiceClient
}

func NewLOMSClient(client inventorypb.InventoryServiceClient) *LOMSClient {
	return &LOMSClient{client: client}
}

func (c *LOMSClient) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	resp, err := c.client.GetStock(ctx, &inventorypb.GetStockRequest{Sku: sku})
	if err != nil {
		return 0, err
	}

	return uint64(resp.Count), nil
}

func (c *LOMSClient) CreateOrder(ctx context.Context, userID int64, items []entity.ListCartItem) (int64, error) {
	pbItems := make([]*inventorypb.CreateOrderItem, 0, len(items))
	for _, item := range items {
		pbItems = append(pbItems, &inventorypb.CreateOrderItem{
			Sku:   item.SKU,
			Count: item.Count,
		})
	}

	resp, err := c.client.CreateOrder(ctx, &inventorypb.CreateOrderRequest{
		UserId: userID,
		Item:   pbItems,
	})
	if err != nil {
		return 0, err
	}

	return resp.OrderId, nil
}

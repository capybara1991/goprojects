package adapter

import (
	"context"

	inventorypb "grpc_pr/gen/inventory"
	"grpc_pr/internal/cart/entity"
)

type ProductClient struct {
	client inventorypb.InventoryServiceClient
}

func NewProductClient(client inventorypb.InventoryServiceClient) *ProductClient {
	return &ProductClient{client: client}
}

func (c *ProductClient) GetProductInfo(ctx context.Context, sku uint32) (entity.ProductInfo, error) {
	resp, err := c.client.GetProductInfo(ctx, &inventorypb.GetProductInfoRequest{Sku: sku})
	if err != nil {
		return entity.ProductInfo{}, entity.ErrProductNotFound
	}

	return entity.ProductInfo{Name: resp.Name, Price: resp.Price}, nil
}

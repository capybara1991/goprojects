package service

import (
	"context"
	"errors"
	"fmt"

	"grpc_pr/internal/cart/entity"
	"grpc_pr/internal/cart/repository"
)

type ProductClient interface {
	GetProductInfo(ctx context.Context, sku uint32) (entity.ProductInfo, error)
}

type LOMSClient interface {
	GetStocks(ctx context.Context, sku uint32) (uint64, error)
	CreateOrder(ctx context.Context, userID int64, items []entity.ListCartItem) (int64, error)
}

type ItemService struct {
	repo          repository.CartRepository
	productClient ProductClient
	lomsClient    LOMSClient
}

func NewItemService(repo repository.CartRepository, productClient ProductClient, lomsClient LOMSClient) *ItemService {
	return &ItemService{
		repo:          repo,
		productClient: productClient,
		lomsClient:    lomsClient,
	}
}

func (s *ItemService) AddItem(ctx context.Context, userID int64, sku, count uint32) error {
	if _, err := s.productClient.GetProductInfo(ctx, sku); err != nil {
		if errors.Is(err, entity.ErrProductNotFound) {
			return fmt.Errorf("get product info: %w", entity.ErrProductNotFound)
		}
		return err
	}

	available, err := s.lomsClient.GetStocks(ctx, sku)
	if err != nil {
		return err
	}

	if uint64(count) > available {
		return fmt.Errorf("requested %d, available %d: %w", count, available, entity.ErrInsufficientStock)
	}

	return s.repo.AddItem(ctx, userID, sku, count)
}

func (s *ItemService) DeleteItem(ctx context.Context, userID int64, sku uint32) error {
	return s.repo.DeleteItem(ctx, userID, sku)
}

func (s *ItemService) ListCart(ctx context.Context, userID int64) ([]entity.ListCartItem, uint64, error) {
	itemsCh, errCh := s.repo.ListItems(ctx, userID)

	var result []entity.ListCartItem
	var totalPrice uint64

	for itemsCh != nil || errCh != nil {
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()

		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			if err != nil {
				return nil, 0, err
			}

		case item, ok := <-itemsCh:
			if !ok {
				itemsCh = nil
				continue
			}

			product, err := s.productClient.GetProductInfo(ctx, item.SKU)
			if err != nil {
				return nil, 0, entity.ErrProductNotFound
			}

			result = append(result, entity.ListCartItem{
				SKU:   item.SKU,
				Count: item.Count,
				Name:  product.Name,
				Price: product.Price,
			})

			totalPrice += uint64(item.Count) * product.Price
		}
	}

	if len(result) == 0 {
		return nil, 0, entity.ErrCartIsEmpty
	}

	return result, totalPrice, nil
}

func (s *ItemService) Checkout(ctx context.Context, userID int64) (int64, error) {
	items, _, err := s.ListCart(ctx, userID)
	if err != nil {
		return 0, err
	}

	orderID, err := s.lomsClient.CreateOrder(ctx, userID, items)
	if err != nil {
		return 0, err
	}

	if err := s.repo.ClearCart(ctx, userID); err != nil {
		return 0, err
	}

	return orderID, nil
}

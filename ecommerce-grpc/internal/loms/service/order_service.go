package service

import (
	"context"

	"grpc_pr/internal/loms/entity"
	"grpc_pr/internal/loms/repository"
)

type OrderService struct {
	ordersRepo repository.OrderRepository
	stocksRepo repository.StocksRepository
}

func NewOrderService(ordersRepo repository.OrderRepository, stocksRepo repository.StocksRepository) *OrderService {
	return &OrderService{ordersRepo: ordersRepo, stocksRepo: stocksRepo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, items []entity.OrderItem) (uint32, error) {
	if len(items) == 0 {
		return 0, entity.ErrEmptyOrderItems
	}

	reserved := make([]entity.OrderItem, 0, len(items))
	for _, item := range items {
		if item.SKU == 0 || item.Count == 0 {
			return 0, entity.ErrEmptyOrderItems
		}

		if err := s.stocksRepo.Reserve(ctx, item.SKU, item.Count); err != nil {
			for _, reservedItem := range reserved {
				_ = s.stocksRepo.Unreserve(ctx, reservedItem.SKU, reservedItem.Count)
			}
			return 0, err
		}

		reserved = append(reserved, item)
	}

	order := &entity.Order{
		UserID: userID,
		Items:  append([]entity.OrderItem(nil), items...),
		Status: entity.StatusAwaitingPayment,
	}

	orderID, err := s.ordersRepo.CreateOrder(ctx, order)
	if err != nil {
		for _, item := range reserved {
			_ = s.stocksRepo.Unreserve(ctx, item.SKU, item.Count)
		}
		return 0, err
	}

	return orderID, nil
}

func (s *OrderService) PayOrder(ctx context.Context, orderID uint32) error {
	order, err := s.ordersRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case entity.StatusPaid:
		return entity.ErrOrderAlreadyPaid
	case entity.StatusCancelled:
		return entity.ErrOrderAlreadyCancelled
	case entity.StatusAwaitingPayment:
		order.Status = entity.StatusPaid
		return s.ordersRepo.UpdateOrder(ctx, order)
	default:
		return entity.ErrInvalidOrderStatus
	}
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uint32) error {
	order, err := s.ordersRepo.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	switch order.Status {
	case entity.StatusPaid:
		return entity.ErrOrderAlreadyPaid
	case entity.StatusCancelled:
		return entity.ErrOrderAlreadyCancelled
	case entity.StatusAwaitingPayment:
		for _, item := range order.Items {
			if err := s.stocksRepo.Unreserve(ctx, item.SKU, item.Count); err != nil {
				return err
			}
		}

		order.Status = entity.StatusCancelled
		return s.ordersRepo.UpdateOrder(ctx, order)
	default:
		return entity.ErrInvalidOrderStatus
	}
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uint32) (*entity.Order, error) {
	return s.ordersRepo.GetOrder(ctx, orderID)
}

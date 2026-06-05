package repository

import (
	"context"
	"sync"

	"grpc_pr/internal/loms/entity"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (uint32, error)
	GetOrder(ctx context.Context, orderID uint32) (*entity.Order, error)
	UpdateOrder(ctx context.Context, order *entity.Order) error
}

type InMemoryOrderRepository struct {
	orders map[uint32]*entity.Order
	mx     sync.Mutex
	nextID uint32
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{orders: make(map[uint32]*entity.Order)}
}

func (r *InMemoryOrderRepository) CreateOrder(ctx context.Context, order *entity.Order) (uint32, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.nextID++
	copyOrder := *order
	copyOrder.ID = r.nextID
	r.orders[copyOrder.ID] = &copyOrder

	return copyOrder.ID, nil
}

func (r *InMemoryOrderRepository) GetOrder(ctx context.Context, orderID uint32) (*entity.Order, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	order, ok := r.orders[orderID]
	if !ok {
		return nil, entity.ErrOrderNotFound
	}

	copyOrder := *order
	copyOrder.Items = append([]entity.OrderItem(nil), order.Items...)
	return &copyOrder, nil
}

func (r *InMemoryOrderRepository) UpdateOrder(ctx context.Context, order *entity.Order) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.orders[order.ID]; !ok {
		return entity.ErrOrderNotFound
	}

	copyOrder := *order
	copyOrder.Items = append([]entity.OrderItem(nil), order.Items...)
	r.orders[copyOrder.ID] = &copyOrder
	return nil
}

package repository

import (
	"context"
	"sync"

	"grpc_pr/internal/cart/entity"
)

type CartRepository interface {
	AddItem(ctx context.Context, userID int64, sku, count uint32) error
	DeleteItem(ctx context.Context, userID int64, sku uint32) error
	ListItems(ctx context.Context, userID int64) (<-chan entity.CartItem, <-chan error)
	ClearCart(ctx context.Context, userID int64) error
}

type InMemoryCartRepository struct {
	mx    sync.RWMutex
	carts map[int64][]entity.CartItem
}

func NewInMemoryCartRepository() *InMemoryCartRepository {
	return &InMemoryCartRepository{carts: make(map[int64][]entity.CartItem)}
}

func (r *InMemoryCartRepository) AddItem(ctx context.Context, userID int64, sku, count uint32) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	items := r.carts[userID]
	for i := range items {
		if items[i].SKU == sku {
			items[i].Count += count
			r.carts[userID] = items
			return nil
		}
	}

	r.carts[userID] = append(items, entity.CartItem{SKU: sku, Count: count})
	return nil
}

func (r *InMemoryCartRepository) DeleteItem(ctx context.Context, userID int64, sku uint32) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	items, ok := r.carts[userID]
	if !ok || len(items) == 0 {
		return entity.ErrCartNotFound
	}

	rebuild := make([]entity.CartItem, 0, len(items))
	found := false
	for _, item := range items {
		if item.SKU == sku {
			found = true
			continue
		}
		rebuild = append(rebuild, item)
	}

	if !found {
		return entity.ErrCartItemNotFound
	}

	if len(rebuild) == 0 {
		delete(r.carts, userID)
		return nil
	}

	r.carts[userID] = rebuild
	return nil
}

func (r *InMemoryCartRepository) ListItems(ctx context.Context, userID int64) (<-chan entity.CartItem, <-chan error) {
	out := make(chan entity.CartItem)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		r.mx.RLock()
		items, ok := r.carts[userID]
		if ok {
			items = append([]entity.CartItem(nil), items...)
		}
		r.mx.RUnlock()

		if !ok || len(items) == 0 {
			errCh <- entity.ErrCartIsEmpty
			return
		}

		for _, item := range items {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case out <- item:
			}
		}
	}()

	return out, errCh
}

func (r *InMemoryCartRepository) ClearCart(ctx context.Context, userID int64) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.carts[userID]; !ok {
		return entity.ErrCartIsEmpty
	}

	delete(r.carts, userID)
	return nil
}

package repository

import (
	"context"
	"sync"

	"grpc_pr/internal/loms/entity"
)

type StocksRepository interface {
	GetStocks(ctx context.Context, sku uint32) (uint64, error)
	SetStocks(ctx context.Context, sku uint32, count uint64) error
	Reserve(ctx context.Context, sku uint32, count uint32) error
	Unreserve(ctx context.Context, sku uint32, count uint32) error
}

type InMemoryStocksRepository struct {
	stocks map[uint32]uint64
	mx     sync.Mutex
}

func NewInMemoryStocksRepository() *InMemoryStocksRepository {
	return &InMemoryStocksRepository{
		stocks: map[uint32]uint64{
			1:  100,
			2:  80,
			3:  45,
			4:  60,
			5:  120,
			6:  70,
			7:  35,
			8:  200,
			12: 100,
			18: 2,
			19: 50,
		},
	}
}

func (r *InMemoryStocksRepository) GetStocks(ctx context.Context, sku uint32) (uint64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	count, ok := r.stocks[sku]
	if !ok {
		return 0, entity.ErrProductNotFound
	}

	return count, nil
}

func (r *InMemoryStocksRepository) SetStocks(ctx context.Context, sku uint32, count uint64) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.stocks[sku] = count
	return nil
}

func (r *InMemoryStocksRepository) Reserve(ctx context.Context, sku uint32, count uint32) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	available, ok := r.stocks[sku]
	if !ok {
		return entity.ErrProductNotFound
	}
	if available < uint64(count) {
		return entity.ErrInsufficientStock
	}

	r.stocks[sku] -= uint64(count)
	return nil
}

func (r *InMemoryStocksRepository) Unreserve(ctx context.Context, sku uint32, count uint32) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.stocks[sku] += uint64(count)
	return nil
}

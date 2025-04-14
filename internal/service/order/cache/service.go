package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/models"
)

const (
	limit = 10
)

type (
	Repository interface {
		GetAll(ctx context.Context, limit int64) ([]models.Order, error)
	}

	Cache interface {
		Put(key int64, value models.Order) error
	}
)

type CacheService struct {
	repo  Repository
	cache Cache
}

func NewCacheService(ctx context.Context, storage Repository, cache Cache, flag bool) (*CacheService, error) {
	orders, err := storage.GetAll(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("storage.ListOrders: %w", err)
	}
	for _, val := range orders {
		if flag && val.Status == consts.ReturnedSt {
			continue
		}
		if err := cache.Put(val.OrderID, val); err != nil {
			log.Printf("cache.AddCache: %v", err)
		}
	}

	return &CacheService{repo: storage, cache: cache}, nil
}

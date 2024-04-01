package store

import (
	"context"
	"errors"
	"time"

	"github.com/JavaHutt/hashcash/configs"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewRedisStore is a constructor
func NewRedisStore(cfg configs.Config) *Store {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.StoreAddr,
		Password: "",
		DB:       0,
	})

	return &Store{
		rdb: rdb,
		ttl: cfg.StoreExpiration,
	}
}

func (s *Store) Set(ctx context.Context, key string) error {
	return s.rdb.Set(ctx, key, "", s.ttl).Err()
}

func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	result, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}

		return false, err
	}

	return result == 1, nil
}

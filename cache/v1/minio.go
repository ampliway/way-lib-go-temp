package v1

import (
	"context"
	"fmt"
	"time"

	"github.com/ampliway/way-lib-go/cache"
	"github.com/ampliway/way-lib-go/helper/reflection"
	"github.com/redis/go-redis/v9"
)

var _ cache.V1 = (*Redis)(nil)

type Redis struct {
	client *redis.Client
	prefix string
}

func New(cfg *Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.CacheEndpoint,
		Password: cfg.CachePassword,
		DB:       0,
	})

	prefix := reflection.AppNamePkg()

	return &Redis{
		prefix: prefix,
		client: client,
	}, nil
}

func (r *Redis) Set(key string, data string, expiration time.Duration) error {
	err := r.client.Set(context.Background(), fmt.Sprintf("%s-%s", r.prefix, key), data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) Get(key string) (string, error) {
	value, err := r.client.Get(context.Background(), fmt.Sprintf("%s-%s", r.prefix, key)).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return value, nil
}

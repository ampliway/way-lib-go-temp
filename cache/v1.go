package cache

import "time"

const MODULE_NAME = "cache"

type V1 interface {
	Set(key string, data string, expiration time.Duration) error
	Get(key string) (string, error)
}

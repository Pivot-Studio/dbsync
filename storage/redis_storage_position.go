package storage

import "github.com/go-redis/redis/v8"

type redisPositionStorage struct {
	rdb *redis.Client
}

package cache

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func RedisConect() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
	return rdb
}

type RedisSetStruct struct {
	key   string
	value string
	ttl   time.Duration
}

func RedisSet(rdb *redis.Client, s RedisSetStruct) error {
	err := rdb.Set(ctx, s.key, s.value, s.ttl).Err()

	if err != nil {
		fmt.Printf("%s -> %s", CannotSetKeyError, err)
		return CannotSetKeyError
	}

	return nil
}

func RedisGet(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()

	if err != nil {
		fmt.Printf("%s -> %s", CannotGetKeyError, err)
		return "", CannotGetKeyError
	}

	return val, nil
}

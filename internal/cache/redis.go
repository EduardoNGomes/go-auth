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
	Key      string
	Value    string
	Duration time.Duration
}

func RedisSet(rdb *redis.Client, s RedisSetStruct) error {
	err := rdb.Set(ctx, s.Key, s.Value, s.Duration).Err()

	if err != nil {
		fmt.Printf("%s -> %s", CannotSetKeyError, err)
		return CannotSetKeyError
	}

	return nil
}

func RedisGetDel(rdb *redis.Client, key string) (string, error) {
	val, err := rdb.GetDel(ctx, key).Result()

	if err != nil {
		fmt.Printf("%s -> %s", CannotGetKeyError, err)
		return "", CannotGetKeyError
	}

	return val, nil
}

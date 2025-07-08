package config

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func NewRedisClient() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", Config.RedisHost, Config.RedisPort)
	dbNum, err := strconv.Atoi(Config.RedisDb)
	if err != nil {
		fmt.Printf("invalid REDIS_DB %q: %v\n", Config.RedisDb, err)
		return nil, err
	}
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: Config.RedisPass,
		DB:       dbNum,
	})

	return client, nil
}

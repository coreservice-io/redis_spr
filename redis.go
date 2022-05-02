package redis_spr

import (
	"context"
	"crypto/tls"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func initRedisClient(addr string, port int, userName string, password string, useTls bool) (redisClient *redis.ClusterClient, err error) {
	config := &redis.ClusterOptions{
		Addrs:    []string{addr + ":" + strconv.Itoa(port)},
		Username: userName,
		Password: password,
	}
	if useTls {
		config.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	redisClient = redis.NewClusterClient(config)

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return
}

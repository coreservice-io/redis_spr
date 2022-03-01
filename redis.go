package RedisSpr

import (
	"context"
	"crypto/tls"
	"strconv"

	"github.com/go-redis/redis/v8"
)

func initRedisClient(addr string, port int, userName string, password string) (redisClient *redis.ClusterClient, err error) {
	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr + ":" + strconv.Itoa(port)},
		Username: userName,
		Password: password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return
}

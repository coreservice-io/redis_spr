package USpr

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
)

func initRedisClient(addr string, port int, userName string, password string) (redisClient *redis.ClusterClient, err error) {
	redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{addr + ":" + strconv.Itoa(port)},
		Username: userName,
		Password: password,
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return
}

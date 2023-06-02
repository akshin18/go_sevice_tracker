package go_server

import (
	"os"

	"github.com/go-redis/redis"
)

func GetRedis() *redis.ClusterClient {

	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     os.Getenv("REDIS_IP_PORT"),
	// 	Password: os.Getenv("REDIS_PWD"),
	// 	DB:       0,
	// })
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{os.Getenv("REDIS_IP_PORT")},
		Password: os.Getenv("REDIS_PWD"),

		// To route commands by latency or randomly, enable one of the following.
		RouteByLatency: true,
		//RouteRandomly: true,
	})

	return rdb
}

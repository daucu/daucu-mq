package daucu_mq

import (
	"crypto/tls"
	"github.com/go-redis/redis/v8"
)


type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int, tlsConfig *tls.Config) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  password,
		DB:        db,
		TLSConfig: tlsConfig, // TLS for encryption
	})

	return &RedisClient{Client: rdb}
}

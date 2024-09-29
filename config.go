package daucu_mq

import (
	"time"
)

func DefaultConfig() Config {
	return Config{
		RedisAddr:     "localhost:6379",
		RedisPassword: "",
		RedisDB:       0,
		MaxRetries:    5,
		VisibilityTimeout: 30 * time.Second,
		TLSConfig:     nil, // Set your TLS config here for security
	}
}

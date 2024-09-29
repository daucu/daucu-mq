package daucu_mq

import (
	"time"
	"crypto/tls"
)

type Message struct {
	ID         string    `json:"id"`
	Data       string    `json:"data"`
	RetryCount int       `json:"retry_count"`
	Timestamp  time.Time `json:"timestamp"`
}

type Config struct {
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	MaxRetries       int
	VisibilityTimeout time.Duration
	TLSConfig        *tls.Config // Optional TLS config for Redis security
}

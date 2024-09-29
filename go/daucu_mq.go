package daucu_mq

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type MessageQueue struct {
	client *redis.Client
	config Config
}

// NewQueue initializes a new MessageQueue with default Redis configuration
func NewQueue(config Config) (*MessageQueue, error) {
	return NewQueueWithCustomRedis(config, &redis.Options{
		Addr:      config.RedisAddr,
		Password:  config.RedisPassword,
		DB:        config.RedisDB,
		TLSConfig: config.TLSConfig, // Optional TLS config for security
	})
}

// NewQueueWithCustomRedis initializes a new MessageQueue with a custom Redis configuration
func NewQueueWithCustomRedis(config Config, redisOptions *redis.Options) (*MessageQueue, error) {
	client := redis.NewClient(redisOptions)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		client: client,
		config: config,
	}, nil
}

// Push adds a new message to the queue
func (q *MessageQueue) Push(queueName string, msg Message) error {
	msgData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = q.client.LPush(ctx, queueName, msgData).Err()
	if err != nil {
		return err
	}

	return nil
}

// Pull retrieves a message from the queue with visibility timeout handling
func (q *MessageQueue) Pull(queueName string) (*Message, error) {
	result := q.client.RPopLPush(ctx, queueName, queueName+":processing").Val()
	if result == "" {
		return nil, redis.Nil
	}

	var msg Message
	if err := json.Unmarshal([]byte(result), &msg); err != nil {
		return nil, err
	}

	q.client.Expire(ctx, queueName+":processing:"+msg.ID, q.config.VisibilityTimeout)

	return &msg, nil
}

// Ack acknowledges successful processing and removes the message from Redis
func (q *MessageQueue) Ack(queueName, msgID string) error {
	return q.client.LRem(ctx, queueName+":processing", 1, msgID).Err()
}

// Retry moves a failed message back to the queue with exponential backoff or DLQ if max retries exceeded
func (q *MessageQueue) Retry(queueName string, msg Message) error {
	if msg.RetryCount >= q.config.MaxRetries {
		fmt.Printf("Message %s sent to Dead Letter Queue\n", msg.ID)
		return q.Push(queueName+":dlq", msg) // Push to Dead Letter Queue
	}

	// Exponential backoff
	backoffDuration := time.Duration(math.Pow(2, float64(msg.RetryCount))) * time.Second
	time.Sleep(backoffDuration)

	msg.RetryCount++
	return q.Push(queueName, msg)
}

// GracefulShutdown waits for all in-progress messages to be processed
func (q *MessageQueue) GracefulShutdown(queueName string) {
	for {
		processingCount, err := q.client.LLen(ctx, queueName+":processing").Result()
		if err != nil || processingCount == 0 {
			break
		}
		fmt.Printf("Waiting for %d messages to be processed\n", processingCount)
		time.Sleep(1 * time.Second)
	}
}

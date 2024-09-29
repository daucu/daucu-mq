package main

import (
	"fmt"
	"time"
	"daucu-mq"

	"github.com/go-redis/redis/v8"
)

func main() {
	config := daucu_mq.DefaultConfig()

	// Example 1: Using Default Redis Configuration
	queue, err := daucu_mq.NewQueue(config)
	if err != nil {
		panic(err)
	}

	// Example 2: Using Custom Redis Configuration
	customRedisOptions := &redis.Options{
		Addr:     "custom-redis-server:6379",
		Password: "custompassword", // no password set by default
		DB:       1,                // use custom DB
	}

	customQueue, err := daucu_mq.NewQueueWithCustomRedis(config, customRedisOptions)
	if err != nil {
		panic(err)
	}

	// Example message
	msg := daucu_mq.Message{
		ID:        "123",
		Data:      "Hello, World!",
		Timestamp: time.Now(),
	}

	// Push message using the custom queue
	if err := customQueue.Push("my_queue", msg); err != nil {
		fmt.Println("Error pushing message:", err)
	}

	// Pull message and process
	pulledMsg, err := customQueue.Pull("my_queue")
	if err != nil {
		fmt.Println("Error pulling message:", err)
	} else {
		fmt.Println("Processing message:", pulledMsg.Data)
		
		// Simulate success or failure
		success := true

		if success {
			customQueue.Ack("my_queue", pulledMsg.ID)
		} else {
			customQueue.Retry("my_queue", *pulledMsg)
		}
	}

	// Graceful shutdown example
	customQueue.GracefulShutdown("my_queue")
}

package main

import (
    "fmt"
    "github.com/daucu/daucu-mq"
)

func main() {

	// Custom Redis configuration
    customConfig := daucu_mq.Config{
        RedisOptions: &redis.Options{
            Addr:     "custom-redis-url:6379", // Redis server address
            Password: "your-password",         // Redis password
            DB:       0,                       // Redis database number
        },
    }

    // Initialize the queue with custom configuration
    queue, err := daucu_mq.NewQueue(customConfig)
    if err != nil {
        panic(err)
    }

    // Initialize the queue with default configuration
    // queue, err := daucu_mq.NewQueue(daucu_mq.DefaultConfig())
    // if err != nil {
    //     panic(err)
    // }

    // Create a new message
    msg := daucu_mq.Message{
        ID:   "1",
        Data: "Hello from daucu-mq",
    }

    // Push message to the queue
    err = queue.Push("test_queue", msg)
    if err != nil {
        fmt.Println("Error pushing message:", err)
    }

    // Pull message from the queue
    pulledMsg, err := queue.Pull("test_queue")
    if err != nil {
        fmt.Println("Error pulling message:", err)
    } else {
        fmt.Printf("Pulled Message: ID=%s, Data=%s\n", pulledMsg.ID, pulledMsg.Data)
    }
}

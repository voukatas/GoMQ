# GoMQ A Simple In-Memory Message Queue in Go

## Description

This project is a basic implementation of an in-memory message queue in Go. It provides simple `Enqueue` and `Dequeue` operations with optional retry mechanisms. The queue is designed to be thread-safe for use in concurrent environments.

## Installation

To install this package, run the following command:

```bash
go get github.com/voukatas/GoMQ
```
# Usage
Here's a simple example that demonstrates how to use the message queue:

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yourusername/my-message-queue/pkg/queue"
)

func main() {
	mq := queue.NewMessageQueue(10)

	// Enqueue a message
	err := mq.Enqueue(queue.Message{ID: "1"}, 3, 10*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to enqueue message: %v", err)
	}

	// Dequeue a message
	msg, err := mq.Dequeue(3, 10*time.Millisecond)
	if err != nil {
		log.Fatalf("Failed to dequeue message: %v", err)
	}

	fmt.Printf("Dequeued message: %s\n", msg.ID)
}

```

# Running Tests
To run tests, navigate to the project directory and execute:
```bash
go test ./...

```



package main

import (
	"fmt"
	"log"
	"time"

	"github.com/voukatas/GoMQ/pkg/queue"
)

func main() {
	mq := queue.NewMessageQueue(10)

	// producers
	for p := 0; p < 3; p++ {
		go func(producerID int) {
			for i := 0; i < 10; i++ {
				err := mq.Enqueue(queue.Message{
					ID:        fmt.Sprintf("P%d-%d", producerID, i),
					Content:   "Hello World!",
					Timestamp: time.Now(),
				}, 3, time.Millisecond*50)
				if err != nil {
					log.Printf("Producer %d: Failed to enqueue message: %v", producerID, err)
				}
				time.Sleep(time.Millisecond * 100)
			}
		}(p)
	}

	// Multiple consumers
	for c := 0; c < 2; c++ {
		go func(consumerID int) {
			for i := 0; i < 15; i++ {
				msg, err := mq.Dequeue(3, time.Millisecond*50)
				if err != nil {
					log.Printf("Consumer %d: Failed to dequeue message: %v", consumerID, err)
				} else {
					log.Printf("Consumer %d: Dequeued message: %s", consumerID, msg.ID)
				}
				time.Sleep(time.Millisecond * 150)
			}
		}(c)
	}

	// Wait for a while to let producers and consumers do their work
	time.Sleep(time.Second * 10)

	mq.ReportMetrics()
}

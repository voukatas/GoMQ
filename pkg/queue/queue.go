package queue

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
)

type Message struct {
	ID        string
	Content   string
	Timestamp time.Time
}

type MessageQueue struct {
	queue       chan Message
	maxSize     int
	enqueued    int64
	dequeued    int64
	enqueueFail int64
	dequeueFail int64
}

var (
	ErrQueueFull  = errors.New("queue is full")
	ErrQueueEmpty = errors.New("queue is empty")
)

func NewMessageQueue(size int) *MessageQueue {
	return &MessageQueue{
		queue:   make(chan Message, size),
		maxSize: size,
	}
}

func (mq *MessageQueue) Enqueue(msg Message, retries int, delay time.Duration) error {
	for i := 0; i <= retries; i++ {
		select {
		case mq.queue <- msg:
			atomic.AddInt64(&mq.enqueued, 1)
			log.Printf("Message enqueued: %s", msg.ID)
			return nil
		default:
			if i < retries {
				time.Sleep(delay)
				continue
			}
			atomic.AddInt64(&mq.enqueueFail, 1)
			log.Printf("Failed to enqueue message after %d retries: %s", retries, msg.ID)
			return ErrQueueFull
		}
	}
	return ErrQueueFull
}

func (mq *MessageQueue) Dequeue(retries int, delay time.Duration) (Message, error) {
	for i := 0; i <= retries; i++ {
		select {
		case msg := <-mq.queue:
			atomic.AddInt64(&mq.dequeued, 1)
			log.Printf("Message dequeued: %s", msg.ID)
			return msg, nil
		default:
			if i < retries {
				time.Sleep(delay)
				continue
			}
			atomic.AddInt64(&mq.dequeueFail, 1)
			log.Printf("Failed to dequeue message after %d retries", retries)
			return Message{}, ErrQueueEmpty
		}
	}
	return Message{}, ErrQueueEmpty
}

func (mq *MessageQueue) ReportMetrics() {
	log.Printf("Total messages enqueued: %d", atomic.LoadInt64(&mq.enqueued))
	log.Printf("Total messages dequeued: %d", atomic.LoadInt64(&mq.dequeued))
	log.Printf("Total enqueue failures: %d", atomic.LoadInt64(&mq.enqueueFail))
	log.Printf("Total dequeue failures: %d", atomic.LoadInt64(&mq.dequeueFail))
	log.Printf("Current queue size: %d", len(mq.queue))
}


package queue

import (
	"testing"
	"time"
)

func TestEnqueueDequeue(t *testing.T) {
	mq := NewMessageQueue(2) // Queue size is 2

	// Test Enqueue
	err := mq.Enqueue(Message{ID: "1"}, 0, 0)
	if err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}

	// Test Dequeue
	msg, err := mq.Dequeue(0, 0)
	if err != nil {
		t.Errorf("Dequeue failed: %v", err)
	}
	if msg.ID != "1" {
		t.Errorf("Dequeue got wrong message: got %v, want %v", msg.ID, "1")
	}

	// Test Queue Full
	err = mq.Enqueue(Message{ID: "2"}, 0, 0)
	if err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}
	err = mq.Enqueue(Message{ID: "3"}, 0, 0)
	if err != nil {
		t.Errorf("Enqueue failed: %v", err)
	}
	err = mq.Enqueue(Message{ID: "4"}, 0, 0)
	if err != ErrQueueFull {
		t.Errorf("Expected ErrQueueFull, got: %v", err)
	}

	// Test Dequeue
	msg, err = mq.Dequeue(0, 0)
	if err != nil {
		t.Errorf("Dequeue failed: %v", err)
	}
	if msg.ID != "2" {
		t.Errorf("Dequeue got wrong message: got %v, want %v", msg.ID, "1")
	}
	// Test Dequeue
	msg, err = mq.Dequeue(0, 0)
	if err != nil {
		t.Errorf("Dequeue failed: %v", err)
	}
	if msg.ID != "3" {
		t.Errorf("Dequeue got wrong message: got %v, want %v", msg.ID, "1")
	}
	// Test Queue Empty
	mq.Dequeue(0, 0) // This should succeed
	_, err = mq.Dequeue(0, 0)
	if err != ErrQueueEmpty {
		t.Errorf("Expected ErrQueueEmpty, got: %v", err)
	}
}

func TestConcurrentEnqueueDequeue(t *testing.T) {
	mq := NewMessageQueue(100)

	// Concurrently enqueue messages
	go func() {
		for i := 0; i < 50; i++ {
			mq.Enqueue(Message{ID: "E"}, 3, 10*time.Millisecond)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Concurrently dequeue messages
	go func() {
		for i := 0; i < 50; i++ {
			_, err := mq.Dequeue(3, 10*time.Millisecond)
			if err != nil {
				t.Errorf("Concurrent Dequeue failed: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Allow some time for the operations to complete
	time.Sleep(1 * time.Second)

	// The numbers should add up correctly 
	if mq.enqueued != 50 || mq.dequeued != 50 {
		t.Errorf("Metrics do not match: enqueued = %d, dequeued = %d", mq.enqueued, mq.dequeued)
	}
}

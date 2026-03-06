package waiter

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestWaiter_PutAndGet_Success(t *testing.T) {
	r := NewWaiter[int]()
	id := "123"
	ctx := context.Background()

	wait, _ := r.Get(ctx, id)
	go func() {
		ok, err := r.Put(ctx, id, 42)
		if err != nil {
			fmt.Println("Put unexpected error: ", err)
		}
		if !ok {
			fmt.Println("Put should return false")
		}
	}()
	val, err := wait()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestWaiter_Put_NoSuchChannel(t *testing.T) {
	r := NewWaiter[string]()
	ctx := context.Background()
	ok, err := r.Put(ctx, "no-exist", "hello")
	if err != nil {
		t.Fatalf("Put unexpected error: %v", err)
	}
	if ok {
		t.Errorf("Put should return false when channel does not exist")
	}
}

func TestWaiter_Get_Duplicate(t *testing.T) {
	r := NewWaiter[int]()
	id := "xx"
	ctx := context.Background()
	r.Get(ctx, id) // First get should succeed

	wait, _ := r.Get(ctx, id)
	_, err := wait()
	if err == nil || !errors.Is(err, errors.New("waiter already exists: "+id)) && err.Error() != "waiter already exists: "+id {
		t.Errorf("Expected waiter already exists error, got: %v", err)
	}
}

func TestWaiter_Get_ContextTimeout(t *testing.T) {
	r := NewWaiter[string]()
	id := "test-timeout"

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	wait, _ := r.Get(ctx, id)
	start := time.Now()
	val, err := wait()
	duration := time.Since(start)

	if err == nil || err != context.DeadlineExceeded {
		t.Errorf("expected context deadline exceeded, got: %v", err)
	}
	if val != "" {
		t.Errorf("expected zero value, got: %v", val)
	}
	if duration < 45*time.Millisecond {
		t.Errorf("expected to wait at least 50ms, waited only %v", duration)
	}
}

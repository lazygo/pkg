package receiver

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestReceiver_PutAndGet_Success(t *testing.T) {
	r := NewReceiver[int]()
	id := "123"
	ctx := context.Background()

	recv := r.Get(ctx, id)
	go func() {
		ok, err := r.Put(ctx, id, 42)
		if err != nil {
			fmt.Println("Put unexpected error: ", err)
		}
		if !ok {
			fmt.Println("Put should return false")
		}
	}()
	val, err := recv()
	if err != nil {
		t.Fatalf("Get() returned error: %v", err)
	}
	if val != 42 {
		t.Errorf("Expected 42, got %d", val)
	}
}

func TestReceiver_Put_NoSuchChannel(t *testing.T) {
	r := NewReceiver[string]()
	ctx := context.Background()
	ok, err := r.Put(ctx, "no-exist", "hello")
	if err != nil {
		t.Fatalf("Put unexpected error: %v", err)
	}
	if ok {
		t.Errorf("Put should return false when channel does not exist")
	}
}

func TestReceiver_Get_Duplicate(t *testing.T) {
	r := NewReceiver[int]()
	id := "xx"
	ctx := context.Background()
	r.Get(ctx, id) // First get should succeed

	recv2 := r.Get(ctx, id)
	_, err := recv2()
	if err == nil || !errors.Is(err, errors.New("response already exists: "+id)) && err.Error() != "response already exists: "+id {
		t.Errorf("Expected response already exists error, got: %v", err)
	}
}

func TestReceiver_Get_ContextTimeout(t *testing.T) {
	r := NewReceiver[string]()
	id := "test-timeout"

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	recv := r.Get(ctx, id)
	start := time.Now()
	val, err := recv()
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

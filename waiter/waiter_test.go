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

// 竞态回归：等待方超时后条目被清理，迟到的 Put 必须干净返回 false 而不是 panic
func TestWaiter_PutAfterTimeout(t *testing.T) {
	r := NewWaiter[int]()
	id := "late-put"

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	wait, _ := r.Get(ctx, id)
	if _, err := wait(); err != context.DeadlineExceeded {
		t.Fatalf("expected deadline exceeded, got: %v", err)
	}

	// 条目已清理：Put 找不到 id，返回 false；绝不允许 send on closed channel
	ok, err := r.Put(context.Background(), id, 1)
	if err != nil || ok {
		t.Errorf("late Put should return (false, nil), got (%v, %v)", ok, err)
	}

	// 同一 id 超时后可重新注册
	wait2, _ := r.Get(context.Background(), id)
	go func() { r.Put(context.Background(), id, 2) }()
	val, err := wait2()
	if err != nil || val != 2 {
		t.Errorf("re-register Get should receive 2, got (%d, %v)", val, err)
	}
}

// 取消函数：wait 未开始时取消，wait 立即返回错误，且条目被清理、同 id 可复用
func TestWaiter_CancelBeforeWait(t *testing.T) {
	r := NewWaiter[string]()
	id := "cancel-first"
	wait, cancelFn := r.Get(context.Background(), id)
	cancelFn()
	if _, err := wait(); err == nil || err.Error() != "waiter closed" {
		t.Errorf("expected waiter closed after cancel, got: %v", err)
	}
	wait2, _ := r.Get(context.Background(), id)
	go func() { r.Put(context.Background(), id, "ok") }()
	if val, err := wait2(); err != nil || val != "ok" {
		t.Errorf("re-register after cancel failed: (%q, %v)", val, err)
	}
}

// 响应与超时同时就绪时必须优先取到响应，而不是随机丢给 ctx.Err()
func TestWaiter_DataPreferredOnTimeout(t *testing.T) {
	r := NewWaiter[int]()
	id := "race-prefer-data"
	ctx, cancel := context.WithCancel(context.Background())
	wait, _ := r.Get(ctx, id)
	if ok, err := r.Put(context.Background(), id, 42); !ok || err != nil {
		t.Fatalf("Put failed: (%v, %v)", ok, err)
	}
	cancel() // 数据已送达后再超时
	val, err := wait()
	if err != nil || val != 42 {
		t.Errorf("delivered data must be preferred over ctx timeout, got (%d, %v)", val, err)
	}
}

// 并发压力：注册/投递/超时/取消交错，-race 下不得有数据竞态或 panic
func TestWaiter_ConcurrentStress(t *testing.T) {
	r := NewWaiter[int]()
	const n = 64
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer func() { done <- struct{}{} }()
			id := fmt.Sprintf("stress-%d", i%8)
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(i%20)*time.Millisecond)
			defer cancel()
			wait, cancelFn := r.Get(ctx, id)
			defer cancelFn()
			go func() {
				time.Sleep(time.Duration(i%10) * time.Millisecond)
				r.Put(context.Background(), id, i)
			}()
			// 结果可能是数据、超时或 closed，均合法；这里只验证不 panic、不泄漏 map 条目
			wait()
		}(i)
	}
	for i := 0; i < n; i++ {
		<-done
	}
}

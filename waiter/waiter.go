package waiter

import (
	"context"
	"errors"
	"sync"
)

// Waiter 是一个泛型结构体，用于通过 id 管理不同的请求条目，以便通信。
// 每个条目 = 容量 1 的数据通道 + done 信号；条目终结时只关闭 done、绝不关闭数据通道，
// 从根源上消除"向已关闭 channel 发送"的竞态。
type Waiter[T any] struct {
	ch map[string]*entry[T]
	mu sync.Mutex
}

// entry 单次请求的状态
type entry[T any] struct {
	data chan T        // 容量 1：Put 永不阻塞，并发 Put 先到先得
	done chan struct{} // 关闭表示请求已终结（已送达/已取消/已超时）
	once sync.Once     // 保证 done 只关闭一次
}

// NewWaiter 创建并返回一个新的 Waiter 实例
func NewWaiter[T any]() *Waiter[T] {
	return &Waiter[T]{
		ch: make(map[string]*entry[T]),
	}
}

// Put 向指定 id 投递数据。第一个投递者成功，其余（含已取消/已超时后的迟到投递）返回 false。
// ctx 用于中断被阻塞的投递（数据通道有缓冲，正常情况下不会阻塞）。
func (r *Waiter[T]) Put(ctx context.Context, id string, data T) (bool, error) {
	r.mu.Lock()
	e, ok := r.ch[id]
	r.mu.Unlock()
	if !ok {
		// 没有对应 id 的条目（未注册或已清理），返回 false
		return false, nil
	}
	select {
	case e.data <- data:
		e.once.Do(func() { close(e.done) })
		return true, nil
	case <-e.done:
		return false, nil
	case <-ctx.Done():
		return false, ctx.Err()
	}
}

// Get 为指定 id 创建请求条目并返回等待函数与取消函数。
// 若该 id 已存在，返回一个立即报错的等待函数。
func (r *Waiter[T]) Get(ctx context.Context, id string) (func() (T, error), func()) {
	var zero T
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.ch[id]; ok {
		// 如果该 id 已存在，返回错误的接收函数
		return func() (T, error) {
			return zero, errors.New("waiter already exists: " + id)
		}, func() {}
	}
	e := &entry[T]{
		data: make(chan T, 1),
		done: make(chan struct{}),
	}
	r.ch[id] = e

	// 等待函数：优先取已投递的数据（响应送达与超时同时就绪时不能丢响应），
	// 否则等待 done（送达/取消/超时）或调用方 ctx 结束
	waitFunc := func() (T, error) {
		defer func() {
			e.once.Do(func() { close(e.done) })
			r.mu.Lock()
			if r.ch[id] == e {
				delete(r.ch, id)
			}
			r.mu.Unlock()
		}()
		select {
		case data := <-e.data:
			return data, nil
		default:
		}
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-e.done:
			select {
			case data := <-e.data:
				return data, nil
			default:
				return zero, errors.New("waiter closed")
			}
		}
	}
	// 取消函数：唤醒等待方并清理条目（map 内仍是同一 entry 时才删除，
	// 避免误删已复用该 id 的新条目）
	cancelFunc := func() {
		e.once.Do(func() { close(e.done) })
		r.mu.Lock()
		if r.ch[id] == e {
			delete(r.ch, id)
		}
		r.mu.Unlock()
	}
	return waitFunc, cancelFunc
}

package waiter

import (
	"context"
	"errors"
	"sync"
)

// Waiter 是一个泛型结构体，用于通过 id 管理不同的通道，以便通信
type Waiter[T any] struct {
	ch map[string]chan<- T // 存储 id 与通道的映射关系
	mu sync.Mutex          // 用于确保并发安全
}

// NewWaiter 创建并返回一个新的 Waiter 实例
func NewWaiter[T any]() *Waiter[T] {
	return &Waiter[T]{
		ch: make(map[string]chan<- T),
	}
}

// Put 向指定 id 的通道发送数据，若通道不存在返回 false，不阻塞于锁内
// ctx 支持超时控制
func (r *Waiter[T]) Put(ctx context.Context, id string, data T) (bool, error) {
	r.mu.Lock()
	ch, ok := r.ch[id]
	r.mu.Unlock()
	if !ok {
		// 没有对应 id 的通道，返回 false
		return false, nil
	}
	select {
	case <-ctx.Done():
		// 上下文取消或超时
		return false, ctx.Err()
	case ch <- data:
	}
	return true, nil
}

// Get 为指定 id 创建响应通道并返回一个等待函数
// 若该 id 已存在，返回一个错误函数，避免重复创建
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
	ch := make(chan T)
	r.ch[id] = ch
	// 返回一个接收函数，等待数据或 ctx 超时/取消
	waitFunc := func() (T, error) {
		defer func() {
			r.mu.Lock()
			delete(r.ch, id) // 接收完毕后，从映射中删除通道
			r.mu.Unlock()
			close(ch)
		}()
		select {
		case <-ctx.Done():
			// 上下文取消或超时
			return zero, ctx.Err()
		case data, ok := <-ch:
			if !ok {
				// 通道已关闭
				return zero, errors.New("waiter closed")
			}
			return data, nil
		}
	}
	cancelFunc := func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.ch, id)
		close(ch)
	}
	return waitFunc, cancelFunc
}

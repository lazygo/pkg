package single

import (
	"sync"
)

// 用于保存实例和每个name的初始化控制
type Singletons[R comparable, T any] struct {
	sync.Map
}

func NewSingletons[R comparable, T any]() *Singletons[R, T] {
	return &Singletons[R, T]{}
}

// singletonHolder 用于保存单例实例及其初始化控制（优化版）
type singletonHolder[R comparable, T any] struct {
	once     sync.Once
	refs     sync.Map
	instance T
	err      error
}

func (holder *singletonHolder[R, T]) hasRefs() bool {
	var has bool
	holder.refs.Range(func(key, value any) bool {
		has = true
		return false
	})
	return has
}

// Peek 查看实例是否存在且有引用计数
func (s *Singletons[R, T]) Peek(name string) (T, []R) {
	var zero T
	holderIface, ok := s.Load(name)
	if !ok {
		return zero, nil
	}
	holder := holderIface.(*singletonHolder[R, T])
	if !holder.hasRefs() {
		return zero, nil
	}
	return holder.instance, s.Refs(name)
}

// Refs 获取所有引用
func (s *Singletons[R, T]) Refs(name string) []R {
	holderIface, ok := s.Load(name)
	if !ok {
		return nil
	}
	holder := holderIface.(*singletonHolder[R, T])
	refs := make([]R, 0)
	holder.refs.Range(func(key, value any) bool {
		refs = append(refs, key.(R))
		return true
	})
	return refs
}

// Get 提供类型T的单例（并发安全高效），通过name区分
func (s *Singletons[R, T]) Get(ref R, name string, init func() T) T {
	holder := s.getHolder(name)
	if ref != *new(R) {
		holder.refs.Store(ref, struct{}{})
	}
	holder.once.Do(func() {
		holder.instance = init()
	})
	return holder.instance
}

// GetWithError 提供错误感知的单例获取（并发安全高效）
func (s *Singletons[R, T]) GetWithError(ref R, name string, init func() (T, error)) (T, error) {
	holder := s.getHolder(name)
	if ref != *new(R) {
		holder.refs.Store(ref, struct{}{})
	}
	holder.once.Do(func() {
		instance, err := init()
		if err != nil {
			holder.err = err
			// 初始化失败立即移除，允许下次重试
			s.Delete(name)
			return
		}
		holder.instance = instance
	})
	return holder.instance, holder.err
}

// Put 减少引用计数，归零时清理实例
// cleanup 可选清理函数
func (s *Singletons[R, T]) Put(ref R, name string, cleanup func(instance T) error) error {
	holderIface, ok := s.Load(name)
	if !ok {
		return nil
	}
	holder := holderIface.(*singletonHolder[R, T])
	if ref != *new(R) {
		holder.refs.Delete(ref)
	}

	if !holder.hasRefs() {
		if cleanup != nil {
			cleanup(holder.instance)
		}
		s.CompareAndDelete(name, holder)
	}
	return nil
}

func (s *Singletons[R, T]) getHolder(name string) *singletonHolder[R, T] {
	holderIface, loaded := s.Load(name)
	var holder *singletonHolder[R, T]
	if loaded {
		holder = holderIface.(*singletonHolder[R, T])
	} else {
		holder = &singletonHolder[R, T]{}
		actual, _ := s.LoadOrStore(name, holder)
		holder = actual.(*singletonHolder[R, T])
	}
	return holder
}

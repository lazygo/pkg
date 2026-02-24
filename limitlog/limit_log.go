package limitlog

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	cleanupInterval = 10 * time.Minute // 清理任务的执行频率
	keyExpire       = 30 * time.Minute // 多久不访问就删除
)

var limiter = NewLimiter()

type limiterEntry struct {
	limiter   *rate.Limiter
	lastVisit time.Time
}

type Limiter struct {
	limiters map[string]*limiterEntry
	mu       sync.Mutex
}

// 创建Limiter同时启动后台清理协程
func NewLimiter() *Limiter {
	l := &Limiter{
		limiters: make(map[string]*limiterEntry),
	}
	go l.cleanupLoop()
	return l
}

func (l *Limiter) limiter(key string, interval time.Duration) *rate.Limiter {
	limit := rate.Every(interval)
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, ok := l.limiters[key]
	if !ok {
		lim := rate.NewLimiter(limit, 1)
		entry = &limiterEntry{limiter: lim, lastVisit: now}
		l.limiters[key] = entry
	} else {
		entry.lastVisit = now
	}
	return entry.limiter
}

func (l *Limiter) cleanupLoop() {
	for {
		time.Sleep(cleanupInterval)
		l.mu.Lock()
		now := time.Now()
		for key, entry := range l.limiters {
			if now.Sub(entry.lastVisit) > keyExpire {
				delete(l.limiters, key)
			}
		}
		l.mu.Unlock()
	}
}

func Log(f func(format string, args ...any), key string, interval time.Duration, format string, args ...any) {
	lim := limiter.limiter(key, interval)
	if !lim.Allow() {
		return
	}
	f(format, args...)
}

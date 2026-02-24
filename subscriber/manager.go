package subscriber

import (
	"sync"
)

type SubscriberManager[U comparable, P comparable] struct {
	mu sync.RWMutex
	// key: 产品名称, value: 订阅者集合
	productSubs map[P]map[U]struct{}
	// key: 订阅者, value: 产品订阅顺序，最近的在末尾
	subscriberProducts map[U][]P
	// key: 订阅者, value: product查重用
	subscriberProductsSet map[U]map[P]struct{}
	// 最大订阅额度
	max int
}

// max为最大订阅额度
func NewSubscriberManager[U comparable, P comparable](max int) *SubscriberManager[U, P] {
	return &SubscriberManager[U, P]{
		productSubs:           make(map[P]map[U]struct{}),
		subscriberProducts:    make(map[U][]P),
		subscriberProductsSet: make(map[U]map[P]struct{}),
		max:                   max,
	}
}

// Subscribe 订阅产品
func (sm *SubscriberManager[U, P]) Subscribe(subID U, product P) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// 已经订阅则直接返回
	if _, ok := sm.productSubs[product]; !ok {
		sm.productSubs[product] = make(map[U]struct{})
	}
	if _, ok := sm.productSubs[product][subID]; ok {
		return
	}

	// 初始化订阅者产品集
	if _, ok := sm.subscriberProducts[subID]; !ok {
		sm.subscriberProducts[subID] = []P{}
	}
	if _, ok := sm.subscriberProductsSet[subID]; !ok {
		sm.subscriberProductsSet[subID] = make(map[P]struct{})
	}

	// 如果订阅额度已满，删除最早的
	if len(sm.subscriberProducts[subID]) >= sm.max {
		oldProduct := sm.subscriberProducts[subID][0]
		sm.subscriberProducts[subID] = sm.subscriberProducts[subID][1:]
		delete(sm.subscriberProductsSet[subID], oldProduct)
		// 同步删除 productSubs 里对应关系
		if subs, ok := sm.productSubs[oldProduct]; ok {
			delete(subs, subID)
			if len(subs) == 0 {
				delete(sm.productSubs, oldProduct)
			}
		}
	}

	// 插入新订阅
	sm.subscriberProducts[subID] = append(sm.subscriberProducts[subID], product)
	sm.subscriberProductsSet[subID][product] = struct{}{}
	sm.productSubs[product][subID] = struct{}{}
}

// Unsubscribe 取消订阅产品
func (sm *SubscriberManager[U, P]) Unsubscribe(subID U, product P) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if subs, ok := sm.productSubs[product]; ok {
		delete(subs, subID)
		if len(subs) == 0 {
			delete(sm.productSubs, product)
		}
	}
	// 从订阅者产品列表中移除
	products, ok := sm.subscriberProducts[subID]
	if !ok {
		return
	}

	// 查找并移除
	newProducts := make([]P, 0, len(products))
	for _, p := range products {
		if p != product {
			newProducts = append(newProducts, p)
		}
	}
	if len(newProducts) == 0 {
		delete(sm.subscriberProducts, subID)
		delete(sm.subscriberProductsSet, subID)
	} else {
		sm.subscriberProducts[subID] = newProducts
		if set, ok := sm.subscriberProductsSet[subID]; ok {
			delete(set, product)
		}
	}
}

// ListSubscribers 返回某产品的所有订阅者
func (sm *SubscriberManager[U, P]) ListSubscribers(product P) []U {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	var res []U
	if subs, ok := sm.productSubs[product]; ok {
		for sub := range subs {
			res = append(res, sub)
		}
	}
	return res
}

// ListProducts 列出订阅者所订阅的产品
func (sm *SubscriberManager[U, P]) ListProducts(subID U) []P {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	products, ok := sm.subscriberProducts[subID]
	if !ok {
		return nil
	}
	return append([]P{}, products...)
}

package subscriber

import (
	"reflect"
	"sort"
	"testing"
)

func TestSubscriberManager_SubscribeAndList(t *testing.T) {
	sm := NewSubscriberManager[uint64, string](3)
	subID := uint64(1)

	// 空时ListProducts应为nil
	if prods := sm.ListProducts(subID); prods != nil {
		t.Errorf("Expected nil, got %v", prods)
	}

	// 订阅productA
	sm.Subscribe(subID, "productA")
	got := sm.ListProducts(subID)
	want := []string{"productA"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
	subs := sm.ListSubscribers("productA")
	if len(subs) != 1 || subs[0] != subID {
		t.Errorf("ListSubscribers failed, got %v", subs)
	}

	// 重复订阅不改变
	sm.Subscribe(subID, "productA")
	got2 := sm.ListProducts(subID)
	if !reflect.DeepEqual(got2, want) {
		t.Errorf("Repeated subscribe should not duplicate: got %v", got2)
	}

	// 订阅多个
	sm.Subscribe(subID, "productB")
	sm.Subscribe(subID, "productC")
	want = []string{"productA", "productB", "productC"}
	got = sm.ListProducts(subID)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}
	// 列表子顺序应正确
	if got[0] != "productA" || got[2] != "productC" {
		t.Errorf("Product order incorrect: %v", got)
	}

	// 超出额度 自动移除最早的
	sm.Subscribe(subID, "productD")
	want = []string{"productB", "productC", "productD"}
	got = sm.ListProducts(subID)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Exceeded: want %v, got %v", want, got)
	}
	// 检查productA移除了订阅者
	subsA := sm.ListSubscribers("productA")
	if len(subsA) != 0 {
		t.Errorf("productA subs expect 0, got %v", subsA)
	}

	// 订阅多个id
	sm.Subscribe(2, "productD")
	sm.Subscribe(2, "productE")
	subsD := sm.ListSubscribers("productD")
	sort.Slice(subsD, func(i, j int) bool { return subsD[i] < subsD[j] })
	if !reflect.DeepEqual(subsD, []uint64{1, 2}) {
		t.Errorf("productD should have [1 2], got %v", subsD)
	}
	subsE := sm.ListSubscribers("productE")
	if len(subsE) != 1 || subsE[0] != 2 {
		t.Errorf("productE should have [2], got %v", subsE)
	}
}

func TestSubscriberManager_Unsubscribe(t *testing.T) {
	sm := NewSubscriberManager[uint64, string](2)
	subID := uint64(9)

	// 先订阅
	sm.Subscribe(subID, "productX")
	sm.Subscribe(subID, "productY")
	got := sm.ListProducts(subID)
	want := []string{"productX", "productY"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v, got %v", want, got)
	}

	// 退订不存在的product无影响
	sm.Unsubscribe(subID, "productZ")
	got = sm.ListProducts(subID)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unsubscribe non-exist changed list: %v", got)
	}

	// 退订productX
	sm.Unsubscribe(subID, "productX")
	want = []string{"productY"}
	got = sm.ListProducts(subID)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("after unsubscribe, want %v, got %v", want, got)
	}
	subsX := sm.ListSubscribers("productX")
	if len(subsX) != 0 {
		t.Errorf("productX subscriber should be 0, got %v", subsX)
	}
	// 再退订剩下的productY
	sm.Unsubscribe(subID, "productY")
	got2 := sm.ListProducts(subID)
	if len(got2) != 0 {
		t.Errorf("after all unsubscribed, expected nil/empty, got %v", got2)
	}
}

func TestSubscriberManager_MultipleSubscribers(t *testing.T) {
	sm := NewSubscriberManager[uint64, string](3)
	for i := 100; i < 105; i++ {
		sm.Subscribe(uint64(i), "alpha")
		sm.Subscribe(uint64(i), "beta")
	}
	// alpha应有5个订阅者
	alphaSubs := sm.ListSubscribers("alpha")
	sort.Slice(alphaSubs, func(i, j int) bool { return alphaSubs[i] < alphaSubs[j] })
	expectedIDs := []uint64{100, 101, 102, 103, 104}
	if !reflect.DeepEqual(alphaSubs, expectedIDs) {
		t.Fatalf("want %v, got %v", expectedIDs, alphaSubs)
	}
	// beta应有5个订阅者
	betaSubs := sm.ListSubscribers("beta")
	sort.Slice(betaSubs, func(i, j int) bool { return betaSubs[i] < betaSubs[j] })
	if !reflect.DeepEqual(betaSubs, expectedIDs) {
		t.Fatalf("want %v, got %v", expectedIDs, betaSubs)
	}

	// 某订阅者的产品列表
	products := sm.ListProducts(103)
	if !reflect.DeepEqual(products, []string{"alpha", "beta"}) {
		t.Errorf("sub 103 should have [alpha beta], got %v", products)
	}
}

func TestSubscriberManager_Unsubscribe_Cleanup(t *testing.T) {
	sm := NewSubscriberManager[uint64, string](2)
	sm.Subscribe(1, "a")
	sm.Unsubscribe(1, "a")
	// 退订完再查订阅者列表，map应清理干净
	if prods := sm.subscriberProducts[1]; prods != nil {
		t.Errorf("subscriberProducts not cleaned, got %v", prods)
	}
	if set := sm.subscriberProductsSet[1]; set != nil {
		t.Errorf("subscriberProductsSet not cleaned, got %v", set)
	}
}

func TestSubscriberManager_Boundary(t *testing.T) {
	sm := NewSubscriberManager[uint64, string](1)
	sm.Subscribe(1, "a")
	sm.Subscribe(1, "b")
	// 只允许最新的
	if got := sm.ListProducts(1); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("Should only have b, got %v", got)
	}
	// a的订阅者应为0
	if subs := sm.ListSubscribers("a"); len(subs) != 0 {
		t.Errorf("a subscribers expect 0, got %v", subs)
	}
	// b的订阅者应为1
	if subs := sm.ListSubscribers("b"); !(len(subs) == 1 && subs[0] == 1) {
		t.Errorf("b subscribers wrong, got %v", subs)
	}
}

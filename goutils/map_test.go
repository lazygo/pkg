package goutils

import (
	"slices"
	"sort"
	"testing"
)

func TestMapKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	keys := MapKeys(m)
	sort.Strings(keys)
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Errorf("MapKeys() = %v", keys)
	}

	// nil map 返回空切片
	if got := MapKeys[int, string](nil); got == nil || len(got) != 0 {
		t.Errorf("MapKeys(nil) = %v, want empty slice", got)
	}

	// 空 map
	if got := MapKeys(map[int]bool{}); got == nil || len(got) != 0 {
		t.Errorf("MapKeys(empty) = %v, want empty slice", got)
	}
}

func TestMapValues(t *testing.T) {
	m := map[string]int{"a": 3, "b": 1, "c": 2}
	values := MapValues(m)
	sort.Ints(values)
	if !slices.Equal(values, []int{1, 2, 3}) {
		t.Errorf("MapValues() = %v", values)
	}

	if got := MapValues[int, string](nil); got == nil || len(got) != 0 {
		t.Errorf("MapValues(nil) = %v, want empty slice", got)
	}
}

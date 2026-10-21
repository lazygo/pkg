package goutils

import (
	"slices"
	"testing"
)

type sliceTestUser struct {
	ID   int
	Name string
}

func testUsers() []sliceTestUser {
	return []sliceTestUser{{1, "a"}, {2, "b"}, {1, "c"}}
}

func TestSlice2Map(t *testing.T) {
	users := testUsers()

	// first=true 保留先出现的元素
	first := Slice2Map(users, true, func(u sliceTestUser) int { return u.ID })
	if len(first) != 2 {
		t.Errorf("Slice2Map(first) len = %d, want 2", len(first))
	}
	if first[1].Name != "a" {
		t.Errorf("Slice2Map(first)[1].Name = %q, want %q", first[1].Name, "a")
	}

	// first=false 保留后出现的元素
	last := Slice2Map(users, false, func(u sliceTestUser) int { return u.ID })
	if last[1].Name != "c" {
		t.Errorf("Slice2Map(last)[1].Name = %q, want %q", last[1].Name, "c")
	}

	// 空切片
	empty := Slice2Map([]sliceTestUser{}, true, func(u sliceTestUser) int { return u.ID })
	if len(empty) != 0 {
		t.Errorf("Slice2Map(empty) len = %d, want 0", len(empty))
	}
}

func TestSlice2KV(t *testing.T) {
	users := testUsers()

	first := Slice2KV(users, true, func(u sliceTestUser, index int) (int, string) {
		return u.ID, u.Name
	})
	if first[1] != "a" || first[2] != "b" {
		t.Errorf("Slice2KV(first) = %v", first)
	}

	last := Slice2KV(users, false, func(u sliceTestUser, index int) (int, string) {
		return u.ID, u.Name
	})
	if last[1] != "c" {
		t.Errorf("Slice2KV(last)[1] = %q, want %q", last[1], "c")
	}

	// index 传递正确
	indices := Slice2KV(users, true, func(u sliceTestUser, index int) (int, int) {
		return u.ID, index
	})
	if indices[1] != 0 || indices[2] != 1 {
		t.Errorf("Slice2KV indices = %v", indices)
	}
}

func TestSliceIndex(t *testing.T) {
	got := SliceIndex([]int{1, 2, 3, 4}, func(n int) (int, int) { return n % 2, n })
	if !slices.Equal(got[0], []int{2, 4}) || !slices.Equal(got[1], []int{1, 3}) {
		t.Errorf("SliceIndex() = %v", got)
	}
}

func TestSliceColumns(t *testing.T) {
	got := SliceColumns(testUsers(), func(u sliceTestUser) int { return u.ID })
	if !slices.Equal(got, []int{1, 2, 1}) {
		t.Errorf("SliceColumns() = %v", got)
	}
}

func TestSliceFilter(t *testing.T) {
	got := SliceFilter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
	if !slices.Equal(got, []int{2, 4}) {
		t.Errorf("SliceFilter() = %v", got)
	}

	all := SliceFilter([]int{1, 2}, func(n int) bool { return true })
	if !slices.Equal(all, []int{1, 2}) {
		t.Errorf("SliceFilter(all) = %v", all)
	}
}

func TestSliceUnique(t *testing.T) {
	got := SliceUnique([]string{"a", "b", "a", "c", "b"}, func(s string) string { return s })
	// 保持首次出现的顺序
	if !slices.Equal(got, []string{"a", "b", "c"}) {
		t.Errorf("SliceUnique() = %v", got)
	}
}

func TestSliceJoin(t *testing.T) {
	if got := SliceJoin([]int{1, 2, 3}, ","); got != "1,2,3" {
		t.Errorf("SliceJoin() = %q, want %q", got, "1,2,3")
	}
	if got := SliceJoin([]int{}, ","); got != "" {
		t.Errorf("SliceJoin(empty) = %q, want %q", got, "")
	}
	if got := SliceJoin([]string{"a", "b"}, "-"); got != "a-b" {
		t.Errorf("SliceJoin() = %q, want %q", got, "a-b")
	}
}

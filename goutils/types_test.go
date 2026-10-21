package goutils

import (
	"testing"
)

func TestString(t *testing.T) {
	if got := String("abc"); got != "abc" {
		t.Errorf("String(string) = %q, want %q", got, "abc")
	}
	if got := String([]byte("abc")); got != "abc" {
		t.Errorf("String([]byte) = %q, want %q", got, "abc")
	}
	if got := String(-123); got != "-123" {
		t.Errorf("String(int) = %q, want %q", got, "-123")
	}
	if got := String(int8(-8)); got != "-8" {
		t.Errorf("String(int8) = %q, want %q", got, "-8")
	}
	if got := String(uint64(18446744073709551615)); got != "18446744073709551615" {
		t.Errorf("String(uint64) = %q, want %q", got, "18446744073709551615")
	}
	if got := String(1.5); got != "1.5" {
		t.Errorf("String(float64) = %q, want %q", got, "1.5")
	}
	if got := String(nil); got != "" {
		t.Errorf("String(nil) = %q, want %q", got, "")
	}
	// 不支持的类型返回默认值
	if got := String(true, "def"); got != "def" {
		t.Errorf("String(bool, def) = %q, want %q", got, "def")
	}
	if got := String(nil, "def"); got != "def" {
		t.Errorf("String(nil, def) = %q, want %q", got, "def")
	}
	// 不支持的类型且无默认值
	if got := String(struct{}{}); got != "" {
		t.Errorf("String(struct) = %q, want %q", got, "")
	}
}

func TestInt64(t *testing.T) {
	if got := Int64(int8(-8)); got != -8 {
		t.Errorf("Int64(int8) = %d, want %d", got, -8)
	}
	if got := Int64(uint16(16)); got != 16 {
		t.Errorf("Int64(uint16) = %d, want %d", got, 16)
	}
	if got := Int64(3.9); got != 3 {
		t.Errorf("Int64(float64) = %d, want %d", got, 3)
	}
	if got := Int64("123"); got != 123 {
		t.Errorf("Int64(string) = %d, want %d", got, 123)
	}
	if got := Int64([]byte("123")); got != 123 {
		t.Errorf("Int64([]byte) = %d, want %d", got, 123)
	}
	// 非法输入返回默认值
	if got := Int64("abc", 7); got != 7 {
		t.Errorf("Int64(abc, 7) = %d, want %d", got, 7)
	}
	if got := Int64(nil, 7); got != 7 {
		t.Errorf("Int64(nil, 7) = %d, want %d", got, 7)
	}
	// 非法输入且无默认值
	if got := Int64("abc"); got != 0 {
		t.Errorf("Int64(abc) = %d, want %d", got, 0)
	}
}

func TestUint64(t *testing.T) {
	if got := Uint64(int8(8)); got != 8 {
		t.Errorf("Uint64(int8) = %d, want %d", got, 8)
	}
	if got := Uint64(uint16(16)); got != 16 {
		t.Errorf("Uint64(uint16) = %d, want %d", got, 16)
	}
	if got := Uint64(3.9); got != 3 {
		t.Errorf("Uint64(float64) = %d, want %d", got, 3)
	}
	if got := Uint64("123"); got != 123 {
		t.Errorf("Uint64(string) = %d, want %d", got, 123)
	}
	if got := Uint64([]byte("123")); got != 123 {
		t.Errorf("Uint64([]byte) = %d, want %d", got, 123)
	}
	// 负数字符串为非法输入
	if got := Uint64("-1", 9); got != 9 {
		t.Errorf("Uint64(-1, 9) = %d, want %d", got, 9)
	}
	if got := Uint64("abc"); got != 0 {
		t.Errorf("Uint64(abc) = %d, want %d", got, 0)
	}
}

func TestInt(t *testing.T) {
	if got := Int(int8(-8)); got != -8 {
		t.Errorf("Int(int8) = %d, want %d", got, -8)
	}
	if got := Int(uint16(16)); got != 16 {
		t.Errorf("Int(uint16) = %d, want %d", got, 16)
	}
	if got := Int(3.9); got != 3 {
		t.Errorf("Int(float64) = %d, want %d", got, 3)
	}
	if got := Int("123"); got != 123 {
		t.Errorf("Int(string) = %d, want %d", got, 123)
	}
	if got := Int([]byte("123")); got != 123 {
		t.Errorf("Int([]byte) = %d, want %d", got, 123)
	}
	if got := Int("abc", 7); got != 7 {
		t.Errorf("Int(abc, 7) = %d, want %d", got, 7)
	}
	if got := Int("abc"); got != 0 {
		t.Errorf("Int(abc) = %d, want %d", got, 0)
	}
}

func TestUint(t *testing.T) {
	if got := Uint(int8(8)); got != 8 {
		t.Errorf("Uint(int8) = %d, want %d", got, 8)
	}
	if got := Uint(uint16(16)); got != 16 {
		t.Errorf("Uint(uint16) = %d, want %d", got, 16)
	}
	if got := Uint(3.9); got != 3 {
		t.Errorf("Uint(float64) = %d, want %d", got, 3)
	}
	if got := Uint("123"); got != 123 {
		t.Errorf("Uint(string) = %d, want %d", got, 123)
	}
	if got := Uint([]byte("123")); got != 123 {
		t.Errorf("Uint([]byte) = %d, want %d", got, 123)
	}
	if got := Uint("-1", 9); got != 9 {
		t.Errorf("Uint(-1, 9) = %d, want %d", got, 9)
	}
	if got := Uint("abc"); got != 0 {
		t.Errorf("Uint(abc) = %d, want %d", got, 0)
	}
}

func TestFloat(t *testing.T) {
	if got := Float(int8(-8)); got != -8 {
		t.Errorf("Float(int8) = %v, want %v", got, -8.0)
	}
	if got := Float(uint16(16)); got != 16 {
		t.Errorf("Float(uint16) = %v, want %v", got, 16.0)
	}
	if got := Float(3.5); got != 3.5 {
		t.Errorf("Float(float64) = %v, want %v", got, 3.5)
	}
	if got := Float("1.25"); got != 1.25 {
		t.Errorf("Float(string) = %v, want %v", got, 1.25)
	}
	if got := Float([]byte("1.25")); got != 1.25 {
		t.Errorf("Float([]byte) = %v, want %v", got, 1.25)
	}
	if got := Float("abc", 7.5); got != 7.5 {
		t.Errorf("Float(abc, 7.5) = %v, want %v", got, 7.5)
	}
	if got := Float("abc"); got != 0 {
		t.Errorf("Float(abc) = %v, want %v", got, 0.0)
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"nil", nil, false},
		{"true", true, true},
		{"false", false, false},
		{"int zero", 0, false},
		{"int nonzero", 1, true},
		{"int8 zero", int8(0), false},
		{"int8 nonzero", int8(1), true},
		{"int64 zero", int64(0), false},
		{"int64 nonzero", int64(-1), true},
		{"uint zero", uint(0), false},
		{"uint nonzero", uint(1), true},
		{"uint8 zero", uint8(0), false},
		{"uint8 nonzero", uint8(1), true},
		{"float32 zero", float32(0), false},
		{"float32 nonzero", float32(0.1), true},
		{"float64 zero", float64(0), false},
		{"float64 nonzero", float64(-0.1), true},
		{"empty string", "", false},
		{"nonempty string", "a", true},
		{"empty slice any", []any{}, false},
		{"nonempty slice any", []any{1}, true},
		{"empty slice string", []string{}, false},
		{"nonempty slice string", []string{"a"}, true},
		{"empty map", map[string]int{}, false},
		{"nonempty map", map[string]int{"a": 1}, true},
		{"empty array", [0]int{}, false},
		{"nonempty array", [1]int{}, true},
		{"zero struct", struct{}{}, false},
		{"nonzero struct", struct{ A int }{A: 1}, true},
		{"nil pointer", (*int)(nil), false},
		{"nonnil pointer", new(int), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Bool(tt.value); got != tt.want {
				t.Errorf("Bool(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

// 空字符串视为无值，返回默认值；能解析出数值（即使为0，含首尾空格）则返回转换结果
func TestZeroValueWithDefault(t *testing.T) {
	if got := String("", "def"); got != "def" {
		t.Errorf(`String("", def) = %q, want %q`, got, "def")
	}
	if got := Int64("", 5); got != 5 {
		t.Errorf(`Int64("", 5) = %d, want %d`, got, 5)
	}
	if got := Int64("0", 5); got != 0 {
		t.Errorf(`Int64("0", 5) = %d, want %d`, got, 0)
	}
	if got := Int64("  -0", 5); got != 0 {
		t.Errorf(`Int64("  -0", 5) = %d, want %d`, got, 0)
	}
	if got := Uint64("", 5); got != 5 {
		t.Errorf(`Uint64("", 5) = %d, want %d`, got, 5)
	}
	if got := Uint64("0", 5); got != 0 {
		t.Errorf(`Uint64("0", 5) = %d, want %d`, got, 0)
	}
	if got := Int("", 5); got != 5 {
		t.Errorf(`Int("", 5) = %d, want %d`, got, 5)
	}
	if got := Int("0", 5); got != 0 {
		t.Errorf(`Int("0", 5) = %d, want %d`, got, 0)
	}
	if got := Int("  -0", 5); got != 0 {
		t.Errorf(`Int("  -0", 5) = %d, want %d`, got, 0)
	}
	if got := Uint("", 5); got != 5 {
		t.Errorf(`Uint("", 5) = %d, want %d`, got, 5)
	}
	if got := Uint("0", 5); got != 0 {
		t.Errorf(`Uint("0", 5) = %d, want %d`, got, 0)
	}
	if got := Float("", 5.5); got != 5.5 {
		t.Errorf(`Float("", 5.5) = %v, want %v`, got, 5.5)
	}
	if got := Float("0", 5.5); got != 0 {
		t.Errorf(`Float("0", 5.5) = %v, want %v`, got, 0.0)
	}
	if got := Float("  -0 ", 5.5); got != 0 {
		t.Errorf(`Float("  -0 ", 5.5) = %v, want %v`, got, 0.0)
	}
	// 转换失败仍然使用默认值
	if got := Int64("abc", 5); got != 5 {
		t.Errorf(`Int64("abc", 5) = %d, want %d`, got, 5)
	}
	if got := String(nil, "def"); got != "def" {
		t.Errorf("String(nil, def) = %q, want %q", got, "def")
	}
}

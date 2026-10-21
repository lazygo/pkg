package goutils

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

func TestGenerateOrderNo(t *testing.T) {
	orderNo := GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
	orderNo = GenerateSerialNumber(1000)
	fmt.Println(orderNo)
}

func TestEncodeUint64(t *testing.T) {
	encoded := EncodeUint64(10000)
	if encoded != "3yR" {
		t.Error("enode fail", encoded)
	}
	decoded := DecodeUint64(encoded)
	if decoded != 10000 {
		t.Error("decode fail", decoded)
	}
}

func TestMask(t *testing.T) {
	str := Mask("123", 1, 1)
	if str != "1*3" {
		t.Error("mask fail", str)
	}
	str = Mask("123", 1, 3)
	if str != "123" {
		t.Error("mask fail", str)
	}
	str = Mask("123", 1, 4)
	if str != "123" {
		t.Error("mask fail", str)
	}
	str = Mask("123", 1, 0)
	if str != "1**" {
		t.Error("mask fail", str)
	}
	str = Mask("123", 1, 0)
	if str != "1**" {
		t.Error("mask fail", str)
	}
	str = Mask("18500230535", 2, 3)
	if str != "18******535" {
		t.Error("mask fail", str)
	}
	str = Mask("lzp9421@qq.com", 2, 3)
	if str != "lz******com" {
		t.Error("mask fail", str)
	}
	str = Mask("", 2, 3)
	if str != "" {
		t.Error("mask fail", str)
	}
}

func TestMask_Chinese(t *testing.T) {
	// 全中文字符串
	str := Mask("张三丰", 1, 1)
	if str != "张*丰" {
		t.Errorf("mask fail for Chinese, got: %q, want: %q", str, "张*丰")
	}
	str = Mask("李小龙", 1, 2)
	if str != "李小龙" {
		t.Errorf("mask fail for Chinese (1,2), got: %q, want: %q", str, "李小龙")
	}
	// 超出头尾
	str = Mask("张三丰", 2, 3)
	if str != "张三丰" {
		t.Errorf("mask fail for Chinese, head+tail>len, got: %q, want: %q", str, "张三丰")
	}
	// 空字符串
	str = Mask("", 1, 1)
	if str != "" {
		t.Errorf("mask fail for Chinese, empty string, got: %q, want: %q", str, "")
	}
	// 长度6的中文
	str = Mask("我是一个大好人", 2, 2)
	if str != "我是***好人" {
		t.Errorf("mask fail for Chinese, long, got: %q, want: %q", str, "我是***好人")
	}

	str = Mask("13012345678", 2, 3)
	if str != "13******678" {
		t.Errorf("mask fail for Chinese, long, got: %q, want: %q", str, "13****678")
	}
}

func TestCutRune(t *testing.T) {
	if got := CutRune("张三丰的剑", 3); got != "张三丰" {
		t.Errorf(`CutRune("张三丰的剑", 3) = %q, want %q`, got, "张三丰")
	}
	if got := CutRune("abc", 10); got != "abc" {
		t.Errorf("CutRune overflow = %q, want %q", got, "abc")
	}
	if got := CutRune("abc", 0); got != "" {
		t.Errorf("CutRune(0) = %q, want %q", got, "")
	}
	if got := CutRune("abc", -1); got != "abc" {
		t.Errorf("CutRune(-1) = %q, want %q", got, "abc")
	}
}

func TestRandStr(t *testing.T) {
	s := RandStr(32)
	if len(s) != 32 {
		t.Errorf("RandStr(32) len = %d, want 32", len(s))
	}
	for _, c := range s {
		if !strings.ContainsRune(letters, c) {
			t.Errorf("RandStr contains invalid char %q", c)
		}
	}

	ls := RandLowerStr(16)
	if len(ls) != 16 {
		t.Errorf("RandLowerStr(16) len = %d, want 16", len(ls))
	}
	for _, c := range ls {
		if !strings.ContainsRune(lowers, c) {
			t.Errorf("RandLowerStr contains invalid char %q", c)
		}
	}

	b33 := RandomBase33String(10)
	if len(b33) != 10 {
		t.Errorf("RandomBase33String(10) len = %d, want 10", len(b33))
	}
	for _, c := range b33 {
		if !strings.ContainsRune(base33Alphabets, c) {
			t.Errorf("RandomBase33String contains invalid char %q", c)
		}
	}
}

func TestRandStrConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				_ = RandStr(8)
				_ = GenerateSerialNumber(uint64(j))
			}
		}()
	}
	wg.Wait()
}

func TestEncryptDecryptUID(t *testing.T) {
	for _, id := range []uint64{0, 1, 100, 123456789, 4294967295} {
		s := EncryptUID(id)
		got, err := DecryptUID(s)
		if err != nil {
			t.Errorf("DecryptUID(%q) error: %v", s, err)
			continue
		}
		if got != id {
			t.Errorf("roundtrip %d -> %q -> %d", id, s, got)
		}
	}
}

func TestBase58(t *testing.T) {
	// 空输入不 panic
	if got := Base58Decode(nil); len(got) != 0 {
		t.Errorf("Base58Decode(nil) = %v, want empty", got)
	}
	if got := Base58Decode([]byte("")); len(got) != 0 {
		t.Errorf(`Base58Decode("") = %v, want empty`, got)
	}
	// 非法字符返回 nil
	if got := Base58Decode([]byte("0")); got != nil {
		t.Errorf(`Base58Decode("0") = %v, want nil`, got)
	}

	// 编解码往返
	for _, v := range []uint64{0, 1, 57, 58, 10000, 1 << 32, 1<<63 - 1, 18446744073709551615} {
		if got := DecodeUint64(EncodeUint64(v)); got != v {
			t.Errorf("roundtrip %d -> %q -> %d", v, EncodeUint64(v), got)
		}
	}

	// 超长输入不 panic
	_ = DecodeUint64("zzzzzzzzzzzzzzzzzzzz")
}

func TestToSnakeString(t *testing.T) {
	cases := map[string]string{
		"XxYy":    "xx_yy",
		"XxYY":    "xx_yy",
		"userID":  "user_id",
		"XX":      "xx",
		"a_B":     "a_b",
		"abc":     "abc",
		"field2X": "field2_x",
	}
	for in, want := range cases {
		if got := ToSnakeString(in); got != want {
			t.Errorf("ToSnakeString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToCamelString(t *testing.T) {
	cases := map[string]string{
		"xx_yy":   "XxYy",
		"user_id": "UserId",
		"a":       "A",
	}
	for in, want := range cases {
		if got := ToCamelString(in); got != want {
			t.Errorf("ToCamelString(%q) = %q, want %q", in, got, want)
		}
	}
}

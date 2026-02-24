package feistel

import (
	"fmt"
	"testing"
)

func TestEncryptUID(t *testing.T) {
	// 正确密钥设置（独立参数）
	key0 := uint16(0x8ff0)
	key1 := uint16(0x4435)

	original := uint32(0x12345678)
	fmt.Println("start encrypt")
	encrypted := EncryptUint32(original, key0, key1)
	fmt.Println("start decrypt")
	decrypted := DecryptUint32(encrypted, key0, key1)

	fmt.Printf("Original: 0x%08x\n", original)
	fmt.Printf("Encrypted: 0x%08x\n", encrypted)
	fmt.Printf("Decrypted: 0x%08x\n", decrypted)

	if decrypted != original {
		t.Fatalf("解密失败，预期 0x%08x 实际 0x%08x", original, decrypted)
	}
}

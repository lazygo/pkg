package feistel

func feistelRound(left, right uint16, key uint16) (uint16, uint16) {
	newLeft := right
	shifted := (right << 5) | (right >> 11)
	f := shifted ^ key
	newRight := left ^ f
	return newLeft, newRight
}

func EncryptUint32(id uint32, key0, key1 uint16) uint32 {
	left := uint16(id >> 16)
	right := uint16(id)

	for i := 0; i < 4; i++ {
		left, right = feistelRound(left, right, key0)
		left, right = feistelRound(left, right, key1)
	}

	// 修正点：交换左右部分后返回
	return (uint32(right) << 16) | uint32(left)
}

func DecryptUint32(encrypted uint32, key0, key1 uint16) uint32 {
	left := uint16(encrypted >> 16)
	right := uint16(encrypted)

	for i := 0; i < 4; i++ {
		left, right = feistelRound(left, right, key1)
		left, right = feistelRound(left, right, key0)
	}

	// 修正点：交换左右部分后返回
	return (uint32(right) << 16) | uint32(left)
}

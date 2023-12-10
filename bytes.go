package ddc

func isZeroSlice(slice []byte) bool {
	var s byte
	for i, c := 0, len(slice); i < c; i++ {
		s |= slice[i]
	}
	return s == 0
}

func xor(payload []byte, xor *byte) {
	for i, c := 0, len(payload); i < c; i++ {
		*xor ^= payload[i]
	}
}

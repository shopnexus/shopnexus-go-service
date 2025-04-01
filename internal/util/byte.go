package util

// MergeBitArray merges two byte slices representing bit arrays using bitwise OR.
func MergeBitArray(a, b []byte) []byte {
	lenA, lenB := len(a), len(b)
	maxLen := lenA
	if lenB > maxLen {
		maxLen = lenB
	}

	result := make([]byte, maxLen)

	for i := 0; i < maxLen; i++ {
		var byteA, byteB byte
		if i < lenA {
			byteA = a[i]
		}
		if i < lenB {
			byteB = b[i]
		}
		result[i] = byteA | byteB
	}

	return result
}

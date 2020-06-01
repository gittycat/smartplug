package smartplug

import "encoding/binary"

func encrypt(str string) []byte {
	n := len(str)
	key := byte(171)
	result := make([]byte, n+4)
	binary.BigEndian.PutUint32(result, uint32(n))

	for i, c := range str {
		result[i+4] = key ^ byte(c)
		key = result[i+4]
	}
	return result
}

func decrypt(buf []byte) string {
	key := byte(171)
	result := make([]byte, len(buf))
	for i, b := range buf {
		result[i] = key ^ b
		key = b
	}
	return string(result)
}

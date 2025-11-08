package lockunlock

// 反转字节切片
func ReverseBytes(data []byte) []byte {
	length := len(data)
	reversed := make([]byte, length)
	for i := 0; i < length; i++ {
		reversed[i] = data[length-1-i]
	}
	return reversed
}

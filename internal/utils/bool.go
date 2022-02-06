package utils

// BTOI converts a boolean to an int.
func BTOI(b bool) (res int8) {
	res = 0
	if b {
		res = 1
	}
	return
}

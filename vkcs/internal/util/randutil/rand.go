package randutil

import "math/rand"

// RandomName returns a random string of letters and digits of passed length.
func RandomName(n int) string {
	charSet := []byte("abcdefghijklmnopqrstuvwxyz012346789")
	result := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, charSet[rand.Intn(len(charSet))])
	}
	return string(result)
}

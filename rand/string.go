package rand

import "math/rand"

const safeChars = "23456789ABCDEFGHJKLMNPQRSTWXYZabcdefghijkmnopqrstuvwxyz"

// SafeString returns a string of given length using safe charset.
func SafeString(length int) string {
	return StringFromSet(length, safeChars)
}

// StringFromSet returns a string of given length from supplied charset.
func StringFromSet(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

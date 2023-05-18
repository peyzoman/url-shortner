package main

import (
	"strings"
)

// By removing special characters from base64 we get the following characters
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func Base62Encode(number uint64) string {
	length := len(chars)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)

	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(chars[(number % uint64(length))])
	}

	return encodedBuilder.String()
}

package utils

import (
	"os"
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

func RemoveDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newURL := strings.Replace(url, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}

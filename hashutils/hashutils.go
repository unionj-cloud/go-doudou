package hashutils

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/unionj-cloud/go-doudou/stringutils"
)

// Sha1 return sha1 string
func Sha1(input string) string {
	if input == "" {
		return "adc83b19e793491b1c6ea0fd8b46cd9f32e592fc"
	}
	return fmt.Sprintf("%x", sha1.Sum([]byte(input)))
}

// Secret2Password return password string
func Secret2Password(username, secret string) string {
	if stringutils.IsEmpty(secret) {
		return Sha1(Sha1(secret) + Sha1(username) + Sha1(secret))
	}
	return Sha1(Sha1(secret[:8]) + Sha1(username) + Sha1(secret[8:]))
}

// Base64 returns base64 string
func Base64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

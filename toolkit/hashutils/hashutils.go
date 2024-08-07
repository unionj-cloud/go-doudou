package hashutils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/unionj-cloud/go-doudou/v2/toolkit/stringutils"
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

// encodeHex is borrowed from https://github.com/google/uuid/blob/44b5fee7c49cf3bcdf723f106b36d56ef13ccc88/uuid.go#L200
func encodeHex(dst []byte, uuid [16]byte) {
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
}

// UUIDByString generates uuid hex string from an arbitrary string
func UUIDByString(input string) string {
	sum := sha1.Sum([]byte(input))
	var uuid [16]byte
	copy(uuid[:], sum[:])
	uuid[6] = (uuid[6] & 0x0f) | uint8((5&0xf)<<4)
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
}

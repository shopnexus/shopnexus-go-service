package vnpay

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"time"
)

// formatTime formats time to string in format yyyyMMddHHmmss
func formatTime(t time.Time) string {
	return t.Format("20060102150405")
}

// sign generates a HMAC signature (SHA512) for the given message using the provided key
func sign(message string, key []byte) string {
	sig := hmac.New(sha512.New, key)
	sig.Write([]byte(message))
	return hex.EncodeToString(sig.Sum(nil))
}

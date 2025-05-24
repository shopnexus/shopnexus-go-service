package vnpay

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
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

// buildSortedQuery builds a sorted query string from the input data.
// Note: use "+" instead of " " or %20 for spaces, as VNPAY uses "+" for spaces in their hash calculation.
func buildSortedQuery(inputData map[string]any) string {
	keys := make([]string, 0, len(inputData))
	for k := range inputData {
		keys = append(keys, k)
	}
	fmt.Printf("Keys: %v\n", keys)
	sort.Strings(keys) // To ensure consistent ordering
	fmt.Printf("Sorted keys: %v\n", keys)

	hashData := ""
	for i, k := range keys {
		encoded := url.QueryEscape(k) + "=" + string(inputData[k].(string))
		if i == 0 {
			hashData += encoded
		} else {
			hashData += "&" + encoded
		}
	}

	//! We have to replace the space with + sign because vnpay use + sign :D
	return strings.ReplaceAll(hashData, " ", "+")
}

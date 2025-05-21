package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"net/url"
	"sort"
	"strings"
)

// HashAlgorithm type
type HashAlgorithm string

const (
	SHA256 HashAlgorithm = "sha256"
	SHA512 HashAlgorithm = "sha512"
	// Add more if needed
)

// GlobalConfig represents the configuration for the payment system
type GlobalConfig struct {
	VnpayHost       string
	PaymentEndpoint string
	Endpoints       *struct {
		PaymentEndpoint string
	}
}

// resolveUrlString combines the base and path
func resolveUrlString(base string, path string) string {
	return strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
}

// buildPaymentUrlSearchParams builds query parameters from the map
func buildPaymentUrlSearchParams(data map[string]interface{}) url.Values {
	params := url.Values{}
	keys := make([]string, 0, len(data))

	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := data[key]
		if val != nil {
			params.Add(key, toString(val))
		}
	}

	return params
}

// createPaymentUrl builds the final payment URL
func CreatePaymentUrl(config GlobalConfig, data map[string]interface{}) (*url.URL, error) {
	paymentEndpoint := config.PaymentEndpoint
	if config.Endpoints != nil && config.Endpoints.PaymentEndpoint != "" {
		paymentEndpoint = config.Endpoints.PaymentEndpoint
	}

	redirectUrl, err := url.Parse(resolveUrlString(config.VnpayHost, paymentEndpoint))
	if err != nil {
		return nil, err
	}

	redirectUrl.RawQuery = buildPaymentUrlSearchParams(data).Encode()
	return redirectUrl, nil
}

// calculateSecureHash generates HMAC hex digest
func CalculateSecureHash(secureSecret, data string, hashAlgorithm HashAlgorithm, bufferEncoding string) (string, error) {
	var hashFunc func() hash.Hash
	switch hashAlgorithm {
	case SHA256:
		hashFunc = sha256.New
	case SHA512:
		hashFunc = sha512.New
	default:
		return "", ErrUnsupportedHash
	}

	mac := hmac.New(hashFunc, []byte(secureSecret))
	mac.Write([]byte(data)) // Assuming UTF-8
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// verifySecureHash compares calculated hash with received hash
func VerifySecureHash(secureSecret, data string, hashAlgorithm HashAlgorithm, receivedHash string) (bool, error) {
	calculatedHash, err := CalculateSecureHash(secureSecret, data, hashAlgorithm, "utf-8")
	if err != nil {
		return false, err
	}
	return strings.EqualFold(calculatedHash, receivedHash), nil
}

// Utility to safely convert interface{} to string
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return strings.TrimSpace(strings.ReplaceAll(fmt.Sprintf("%v", v), "\n", ""))
	}
}

// Custom error
var ErrUnsupportedHash = errors.New("unsupported hash algorithm")

package lib

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// HelloRequest represents a request containing a name.
type HelloRequest struct {
	Name string `json:"name"`
}

// StreamConvertRequest represents a request to convert a stream with a URL.
type StreamConvertRequest struct {
	URL string `json:"url"`
}

// StreamConvertResponse represents the response after converting a stream.
type StreamConvertResponse struct {
	URL string `json:"url"`
}

// Config represents the configuration for the application.
type Config struct {
	Port   int
	Secret string
	URL    string
	FFmpeg string
}

// ParseData decodes JSON data into the provided structure.
// It returns either the decoded data or an error.
func ParseData(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// IsParseError checks if the provided error is a parsing error.
func IsParseError(err error) bool {
	return errors.Is(err, &json.SyntaxError{}) || errors.Is(err, &json.UnmarshalTypeError{})
}

// AsyncForEach iterates over a slice and applies the callback function to each element.
// Note: Go handles concurrency differently, and this function runs callbacks synchronously.
// To achieve asynchronous behavior, consider using goroutines.
func AsyncForEach[T any](array []T, callback func(el T, i int)) {
	for i, el := range array {
		callback(el, i)
	}
}

// Echo returns a simple "OK" string.
func Echo() string {
	return "OK"
}

// GetStreamID generates an MD5 hash of the provided URL.
// It returns the hexadecimal representation of the hash.
func GetStreamID(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}

package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ParseResponse unmarshals an HTTP response body into the target struct
func ParseResponse(resp *http.Response, target interface{}) error {
	if resp == nil {
		return fmt.Errorf("nil response")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	// Restore body for potential future reads
	resp.Body = io.NopCloser(io.NopCloser(io.NopCloser(nil)))

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w\nBody: %s", err, string(body))
	}

	return nil
}

// ParseJSON unmarshals JSON bytes into the target struct
func ParseJSON(data []byte, target interface{}) error {
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nil
}

// ReadBody reads and returns the response body as bytes
func ReadBody(resp *http.Response) ([]byte, error) {
	if resp == nil || resp.Body == nil {
		return nil, fmt.Errorf("nil response or body")
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ReadBodyString reads and returns the response body as string
func ReadBodyString(resp *http.Response) (string, error) {
	body, err := ReadBody(resp)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// IsSuccessStatusCode returns true for 2xx status codes
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// IsRedirectStatusCode returns true for 3xx status codes
func IsRedirectStatusCode(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// IsClientErrorStatusCode returns true for 4xx status codes
func IsClientErrorStatusCode(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// IsServerErrorStatusCode returns true for 5xx status codes
func IsServerErrorStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

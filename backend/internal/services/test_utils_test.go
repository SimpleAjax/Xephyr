package services_test

import (
	"github.com/google/uuid"
)

// stringToUUID converts a string to uuid.UUID for test fixtures
// For test IDs that aren't valid UUIDs, generates a deterministic UUID
func stringToUUID(s string) uuid.UUID {
	// Try to parse as UUID first
	if id, err := uuid.Parse(s); err == nil {
		return id
	}
	// Generate deterministic UUID from string using MD5 hash
	return uuid.NewMD5(uuid.NameSpaceOID, []byte(s))
}

// uuidToString converts uuid.UUID to string (helper for test readability)
func uuidToString(u uuid.UUID) string {
	return u.String()
}

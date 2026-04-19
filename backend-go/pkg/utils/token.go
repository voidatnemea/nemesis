package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashToken returns a deterministic SHA-256 hex digest of the provided token.
// It is used to store session tokens in the database without persisting the raw value.
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// GenerateSessionToken produces a user-facing token along with its hashed
// representation for storage. The caller is responsible for persisting the hash
// and returning the raw token to the client.
func GenerateSessionToken() (raw string, hashed string, err error) {
	raw = fmt.Sprintf("fp_authtoken_%s", GenerateRandomString(32))
	if raw == "" {
		return "", "", fmt.Errorf("failed to generate random bytes")
	}
	return raw, HashToken(raw), nil
}

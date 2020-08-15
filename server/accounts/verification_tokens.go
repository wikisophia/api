package accounts

import (
	"crypto/rand"
	"encoding/base64"
)

// Make a new URL-friendly verification token.
func newVerificationToken(length int) (string, error) {
	requiredLength := base64.URLEncoding.DecodedLen(length)
	tmp := make([]byte, requiredLength)
	if _, err := rand.Read(tmp); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tmp), nil
}

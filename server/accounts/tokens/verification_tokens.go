package tokens

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// ResetTokenExpiry determines how long a reset token should be valid for
const ResetTokenExpiry = 48 * time.Hour

// Make a new URL-friendly verification token.
func NewVerificationToken(length int) (string, error) {
	requiredLength := base64.URLEncoding.DecodedLen(length)
	tmp := make([]byte, requiredLength)
	if _, err := rand.Read(tmp); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tmp), nil
}

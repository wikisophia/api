package accounts

import "time"

// Account has the info which is tied to the email which signed up.
type Account struct {
	ID         int64
	Email      string
	ResetToken string
}

// ResetTokenExpiry determines how long a reset token should be valid for
const ResetTokenExpiry = 48 * time.Hour

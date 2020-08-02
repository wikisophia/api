package accounts

import "time"

// Account has the info which is tied to the email which signed up.
type Account struct {
	ID         int64  `json:"-"`
	Email      string `json:"email"`
	ResetToken string `json:"-"`
}

// ResetTokenExpiry determines how long a reset token should be valid for
const ResetTokenExpiry = 48 * time.Hour

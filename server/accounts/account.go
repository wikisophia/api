package accounts

import "time"

// Account is a person who edits arguments on the site.
type Account struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	ResetToken string `json:"-"`
}

// ResetTokenExpiry determines how long a reset token should be valid for
const ResetTokenExpiry = 48 * time.Hour

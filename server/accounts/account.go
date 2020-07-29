package accounts

import "time"

// Account is a person who edits arguments on the site.
type Account struct {
	ID                int64  `json:"id"`
	Email             string `json:"email"`
	VerificationToken string `json:"-"`
}

// A VerificationLink bundles the email and verification token used to activate an account.
// The token must be safe to use in URLs.
type VerificationLink struct {
	Email string
	Token string
}

// VerificationLinkExpiry determines how long a user has to click a
// verificaiton link before it expires.
const VerificationLinkExpiry = 48 * time.Hour

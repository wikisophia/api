package accounts

// Account has the info which is tied to the email which signed up.
type Account struct {
	ID         int64
	Email      string
	ResetToken string
}

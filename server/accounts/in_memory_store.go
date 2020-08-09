package accounts

import "strconv"

// NewMemoryStore makes an empty InMemoryStore with all its variables initialized.
func NewMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		nextID:    1,
		nextReset: 1,
		accounts:  make(map[string]*accountInfo, 1),
	}
}

// InMemoryStore saves accounts in program memory.
// This is mainly intended for testing and easier dev environment setups.
type InMemoryStore struct {
	nextID    int64
	nextReset int64
	accounts  map[string]*accountInfo
}

type accountInfo struct {
	id       int64
	token    string
	password string
}

// NewResetTokenWithAccount sets a password reset token for an account.
// If the email doesn't exist yet, an account will be created for it.
func (s *InMemoryStore) NewResetTokenWithAccount(email string) (string, error) {
	if accountInfo, ok := s.accounts[email]; ok {
		accountInfo.token = "token-" + strconv.FormatInt(accountInfo.id, 10) +
			"-reset-" + strconv.FormatInt(s.nextReset, 10)
		s.nextReset++
		return accountInfo.token, nil
	}
	info := &accountInfo{
		id:       s.nextID,
		token:    "token-" + strconv.FormatInt(s.nextID, 10),
		password: "",
	}
	s.nextID++
	s.accounts[email] = info
	return info.token, nil
}

// SetPassword changes the password associated with the email and returns the account's ID.
// If the email doesn't exist, it returns an EmailNotExistsError.
// If the resetToken is wrong (expired or never returned by ResetPassword(email)),
//   it returns an InvalidPasswordError.
func (s *InMemoryStore) SetPassword(email, password, resetToken string) (int64, error) {
	info, ok := s.accounts[email]
	if !ok {
		return -1, EmailNotExistsError{
			Email: email,
		}
	}
	if info.token == "" || info.token != resetToken {
		return -1, InvalidResetTokenError{}
	}
	if password == "" {
		return -1, ProhibitedPasswordError{}
	}
	info.token = ""
	info.password = password
	return info.id, nil
}

// Authenticate returns the account's ID.
// If the email doesn't exist, it returns an EmailNotExistsError.
// If the password is wrong, it returns an InvalidPasswordError.
func (s *InMemoryStore) Authenticate(email, password string) (int64, error) {
	userInfo, ok := s.accounts[email]
	if !ok {
		return -1, EmailNotExistsError{email}
	}
	if userInfo.password == "" {
		return -1, InvalidPasswordError{}
	}
	if userInfo.password != password {
		return -1, InvalidPasswordError{}
	}
	return userInfo.id, nil
}

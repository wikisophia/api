package accounts

import (
	"context"
)

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
	account  Account
	password string
}

// See the docs on interfaces in store.go
func (s *InMemoryStore) NewResetToken(ctx context.Context, email string) (Account, error) {
	token, err := newVerificationToken(20)
	if err != nil {
		return Account{}, err
	}
	if accountInfo, ok := s.accounts[email]; ok {
		accountInfo.account.ResetToken = token
		s.nextReset++
		return accountInfo.account, nil
	}

	info := &accountInfo{
		account: Account{
			ID:         s.nextID,
			Email:      email,
			ResetToken: token,
		},
		password: "",
	}
	s.nextID++
	s.accounts[email] = info
	return info.account, nil
}

// See the docs on interfaces in store.go
func (s *InMemoryStore) SetForgottenPassword(ctx context.Context, id int64, password, resetToken string) error {
	if password == "" {
		return ProhibitedPasswordError{}
	}
	if resetToken == "" {
		return InvalidResetTokenError{}
	}

	for _, accountInfo := range s.accounts {
		if accountInfo.account.ID == id {
			if resetToken != accountInfo.account.ResetToken {
				return InvalidResetTokenError{}
			}
			accountInfo.password = password
			accountInfo.account.ResetToken = ""
			return nil
		}
	}
	return AccountNotExistsError{}
}

// See the docs on interfaces in store.go
func (s *InMemoryStore) ChangePassword(ctx context.Context, id int64, oldPassword, newPassword string) error {
	if newPassword == "" {
		return ProhibitedPasswordError{}
	}
	for _, accountInfo := range s.accounts {
		if accountInfo.account.ID == id {
			if oldPassword != accountInfo.password {
				return InvalidPasswordError{}
			}
			accountInfo.password = newPassword
			return nil
		}
	}
	return AccountNotExistsError{}
}

// See the docs on interfaces in store.go
func (s *InMemoryStore) Authenticate(ctx context.Context, email, password string) (int64, error) {
	userInfo, ok := s.accounts[email]
	if !ok {
		return -1, AccountNotExistsError{email}
	}
	if userInfo.password == "" {
		return -1, InvalidPasswordError{}
	}
	if userInfo.password != password {
		return -1, InvalidPasswordError{}
	}
	return userInfo.account.ID, nil
}

// See the docs on interfaces in store.go
func (s *InMemoryStore) Close() error {
	return nil
}

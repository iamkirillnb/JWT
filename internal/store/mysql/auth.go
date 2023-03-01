package mysql

import "github.com/uptrace/bun"

type AuthStore struct {
	*bun.DB
}

func NewAuthStore(debug bool) (*AuthStore, error) {
	con, err := setConnection(debug)
	if err != nil {
		return &AuthStore{}, err
	}
	return &AuthStore{con}, nil
}

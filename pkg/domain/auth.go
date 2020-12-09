package domain

import "context"

type (
	Auth struct {
		PlayerID int
		Token    string
	}

	AuthRepository interface {
		Set(ctx context.Context, auth *Auth) error
		Get(ctx context.Context) (*Auth, error)
	}
)

package domain

import "context"

type (
	LoginResponse struct {
		PlayerID int
		Token    string
	}

	LoginRepository interface {
		Login(ctx context.Context, code string) (*LoginResponse, error)
	}
)

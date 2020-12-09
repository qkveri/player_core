package repositories

import (
	"context"
	"database/sql"

	"github.com/qkveri/player_core/pkg/domain"
)

type authSqlLiteRepo struct {
	db *sql.DB
}

func NewAuthSqlLiteRepo(db *sql.DB) *authSqlLiteRepo {
	return &authSqlLiteRepo{
		db: db,
	}
}

func (a *authSqlLiteRepo) Set(ctx context.Context, auth *domain.Auth) error {
	return nil
}

func (a *authSqlLiteRepo) Get(ctx context.Context) (*domain.Auth, error) {
	return nil, nil
}

package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
)

type loginApiRepo struct {
	client api.Client
}

func NewLoginApiRepo(client api.Client) *loginApiRepo {
	return &loginApiRepo{
		client: client,
	}
}

func (l *loginApiRepo) Login(ctx context.Context, code string) (*domain.LoginResponse, error) {
	data := struct {
		Code string `json:"code"`
	}{
		Code: code,
	}

	resRaw, err := l.client.POST(ctx, "/auth/login", data)

	if err != nil {
		return nil, err
	}

	var resLoginResponse struct {
		PlayerID int    `json:"playerId"`
		Token    string `json:"token"`
	}

	if err := json.Unmarshal(resRaw, &resLoginResponse); err != nil {
		return nil, fmt.Errorf("login response unmarshall fail: %w, resRaw: %s", err, resRaw)
	}

	return &domain.LoginResponse{
		PlayerID: resLoginResponse.PlayerID,
		Token:    resLoginResponse.Token,
	}, nil
}

package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
)

type playerInfoApiRepo struct {
	client api.Client
}

func NewPlayerInfoApiRepo(client api.Client) *playerInfoApiRepo {
	return &playerInfoApiRepo{
		client: client,
	}
}

func (p *playerInfoApiRepo) Get(ctx context.Context) (*domain.PlayerInfo, error) {
	resRaw, err := p.client.GET(ctx, "/player/info")

	if err != nil {
		return nil, err
	}

	var resPlayerInfo struct {
		Id           int     `json:"id"`
		Name         string  `json:"name"`
		HasCrossFade bool    `json:"hasCrossfade"`
		CrossFadeSec int     `json:"crossfadeSec"`
		ServerTime   int64   `json:"serverTime"`
		DemoAt       *string `json:"demoAt"`

		Company struct {
			Id      int    `json:"id"`
			Name    string `json:"name"`
			LkURL   string `json:"lkURL"`
			SiteURL string `json:"siteURL"`

			Phone    *string `json:"phone"`
			Email    *string `json:"email"`
			Telegram *string `json:"telegram"`
			Whatsapp *string `json:"whatsapp"`
			Viber    *string `json:"viber"`

			ColorPrimary *string `json:"colorPrimary"`
			LogoLightURL *string `json:"logoLightURL"`
			LogoDarkURL  *string `json:"logoDarkURL"`
		}
	}

	if err := json.Unmarshal(resRaw, &resPlayerInfo); err != nil {
		return nil, fmt.Errorf("player info unmarshall fail: %w", err)
	}

	var demoAt *time.Time

	if resPlayerInfo.DemoAt != nil {
		if da, err := time.Parse(time.RFC3339Nano, *resPlayerInfo.DemoAt); err != nil {
			return nil, fmt.Errorf("player info parse demoAt fail: %w, layout: %s, demoAt: %s", err,
				time.RFC3339Nano,
				*resPlayerInfo.DemoAt)
		} else {
			demoAt = &da
		}
	}

	return &domain.PlayerInfo{
		ID:                resPlayerInfo.Id,
		Name:              resPlayerInfo.Name,
		HasCrossFade:      resPlayerInfo.HasCrossFade,
		CrossFadeDuration: time.Second * time.Duration(resPlayerInfo.CrossFadeSec),
		ServerTime:        time.Unix(resPlayerInfo.ServerTime, 0),
		DemoAt:            demoAt,
		Company: domain.PlayerInfoCompany{
			ID:           resPlayerInfo.Company.Id,
			Name:         resPlayerInfo.Company.Name,
			LkURL:        resPlayerInfo.Company.LkURL,
			SiteURL:      resPlayerInfo.Company.SiteURL,
			Phone:        resPlayerInfo.Company.Phone,
			Email:        resPlayerInfo.Company.Email,
			Telegram:     resPlayerInfo.Company.Telegram,
			Whatsapp:     resPlayerInfo.Company.Whatsapp,
			Viber:        resPlayerInfo.Company.Viber,
			ColorPrimary: resPlayerInfo.Company.ColorPrimary,
			LogoLightURL: resPlayerInfo.Company.LogoLightURL,
			LogoDarkURL:  resPlayerInfo.Company.LogoDarkURL,
		},
	}, nil
}

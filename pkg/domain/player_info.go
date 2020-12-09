package domain

import (
	"context"
	"time"
)

type (
	PlayerInfo struct {
		Id                int
		Name              string
		HasCrossFade      bool
		CrossFadeDuration time.Duration
		ServerTime        time.Time
		DemoAt            *time.Time

		Company PlayerInfoCompany
	}

	PlayerInfoCompany struct {
		Id      int
		Name    string
		LkURL   string
		SiteURL string

		Phone    *string
		Email    *string
		Telegram *string
		Whatsapp *string
		Viber    *string

		ColorPrimary *string
		LogoLightURL *string
		LogoDarkURL  *string
	}

	PlayerInfoRepository interface {
		Get(ctx context.Context) (*PlayerInfo, error)
	}
)

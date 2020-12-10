package app

import (
	"context"
	"strings"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
)

func (a *App) Login(ctx context.Context, callback CallbackLogin, code string) {
	if callback == nil {
		a.logger.Error().Msg("CallbackLogin is nil")
		return
	}

	a.logger.Info().Str("code", code).Msg("login...")

	loginResponse, err := a.loginRepo.Login(ctx, code)

	if err != nil {
		a.logger.Warn().Err(err).Msg("login failed")

		if e, ok := err.(*api.ValidationError); ok {
			if codeErrors, hasCodeError := e.ValidationFails["code"]; hasCodeError {
				callback.SendCodeIncorrectErrorMessage(strings.Join(codeErrors, " "))
				return
			}
		}

		callback.SendErrorMessage(err.Error())

		return
	}

	a.logger.Info().Interface("response", loginResponse).Msg("login success")

	// save auth data to local repo...
	auth := &domain.Auth{
		PlayerID: loginResponse.PlayerID,
		Token:    loginResponse.Token,
	}

	a.logger.Info().Interface("auth", auth).Msg("auth set to repo...")

	if err := a.authRepo.Set(ctx, auth); err != nil {
		a.logger.Err(err).Msg("auth set to repo failed")
		callback.SendErrorMessage(err.Error())

		return
	}

	a.logger.Info().Msg("auth set to repo success")

	a.showScreen(ScreenLoadingData)
}

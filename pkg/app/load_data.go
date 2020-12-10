package app

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qkveri/player_core/pkg/api"
)

func (a *App) LoadData(ctx context.Context, callback CallbackLoadData) {
	if callback == nil {
		a.logger.Error().Msg("CallbackLoadData is nil")
		return
	}

	// load auth data...
	a.logger.Debug().Msg("auth get from repo...")

	auth, err := a.authRepo.Get(ctx)

	if err != nil {
		a.logger.Warn().Err(err).Msg("auth get from repo failed")
	}

	if auth != nil {
		a.logger.Debug().Interface("auth", auth).
			Msg("auth get from repo success, set to api client...")

		a.apiClient.SetAuthToken(auth.Token)

		// inject player id to logger
		a.logger = a.logger.With().Int("playerID", auth.PlayerID).Logger()
	} else {
		a.logger.Debug().Msg("auth not set, show login screen...")
		a.showScreen(ScreenLogin)

		return
	}

	if err := a.loadDataFromAPI(ctx, callback); err != nil {
		a.logger.Warn().Err(err).Msg("load data from API failed")

		switch err.(type) {
		default:
			callback.SendErrorMessage(a.prepareErr(err).Error())

		case *api.UnauthorizedError:
			a.showScreen(ScreenLogin)
		}

		return
	}

	a.showScreen(ScreenPlayer)
}

func (a *App) loadDataFromAPI(ctx context.Context, callback CallbackLoadData) error {
	if err := a.loadPlayerInfo(ctx, callback); err != nil {
		return err
	}

	if err := a.loadMusicData(ctx, callback); err != nil {
		return err
	}

	return nil
}

func (a *App) loadPlayerInfo(ctx context.Context, callback CallbackLoadData) error {
	a.logger.Debug().Msg("playerInfo load starts...")
	callback.SendText("Загрузка информации о заведении...")

	playerInfo, err := a.playerInfoRepo.Get(ctx)

	if err != nil {
		return err
	}

	a.logger.Debug().Interface("playerInfo", playerInfo).Msg("playerInfo loaded")

	// send playerInfo...
	jsonPlayerInfo, err := json.Marshal(playerInfo)

	if err != nil {
		callback.SendErrorMessage(a.prepareErr(err).Error())

		return fmt.Errorf("playerInfo cannot unmarshall: %w", err)
	}

	callback.SendPlayerInfo(string(jsonPlayerInfo))

	return nil
}

func (a *App) loadMusicData(ctx context.Context, callback CallbackLoadData) error {
	a.logger.Debug().Msg("musicData load starts...")
	callback.SendText("Загрузка музыкальных настроек...")

	musicData, err := a.musicDataRepo.Get(ctx)

	if err != nil {
		return err
	}

	a.logger.Debug().Interface("musicData", musicData).Msg("musicData loaded")

	return nil
}

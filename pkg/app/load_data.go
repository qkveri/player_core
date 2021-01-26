package app

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	} else {
		a.logger.Debug().Msg("auth not set, show login screen...")
		a.showScreen(ScreenLogin)

		return
	}

	if err := a.loadDataFromAPI(ctx, callback); err != nil {
		a.logger.Warn().Err(err).Msg("load data from API failed")

		switch err.(type) {
		default:
			callback.SendErrorMessage(a.errMessageForClient(err))

		case *api.UnauthorizedError:
			a.showScreen(ScreenLogin)
		}

		return
	}

	if err := a.awaitLoadFirstTrack(ctx, callback); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		callback.SendErrorMessage(a.errMessageForClient(err))
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

	a.state.PlayerInfo.Lock()
	a.state.PlayerInfo.Set(playerInfo)
	a.state.PlayerInfo.Unlock()

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

	a.state.MusicData.Lock()
	a.state.MusicData.Set(musicData)
	a.state.MusicData.Unlock()

	return nil
}

func (a *App) awaitLoadFirstTrack(ctx context.Context, callback CallbackLoadData) error {
	callback.SendText("Загрузка данных...")

	t := time.NewTicker(time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-t.C:
			a.state.Playlist.RLock()
			progress := a.state.Playlist.FirstItemDownloadProgress()

			a.state.Playlist.RUnlock()

			if progress.IsDone() {
				return nil
			}

			callback.SendText(fmt.Sprintf("Загрузка данных (%s)...", progress))
		}
	}
}

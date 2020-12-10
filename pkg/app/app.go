package app

import (
	"context"
	"path"
	"strings"

	"github.com/rs/zerolog"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/domain/repositories"
)

type App struct {
	config       Config
	callbackMain CallbackMain

	logger    zerolog.Logger
	apiClient api.Client

	// repos...
	playerInfoRepo domain.PlayerInfoRepository
	loginRepo      domain.LoginRepository
	authRepo       domain.AuthRepository
}

func NewApp(config Config, callbackMain CallbackMain) *App {
	a := &App{
		config:       config,
		callbackMain: callbackMain,
	}

	a.logger = a.iniLogger()

	return a
}

func (a *App) Init() {
	// init common...
	a.apiClient = api.NewHTTPClient(a.config.ApiBaseURL)

	// init repos...
	a.playerInfoRepo = repositories.NewPlayerInfoApiRepo(a.apiClient)
	a.loginRepo = repositories.NewLoginApiRepo(a.apiClient)

	authFilePath := path.Join(a.config.DataDir, "a.tk")
	a.authRepo = repositories.NewAuthFileRepo(authFilePath, a.config.SecretKey)

	// show first screen...
	a.showScreen(ScreenLoadingData)
}

func (a *App) LoadingData(ctx context.Context, callback CallbackLoadingData) {
	if callback == nil {
		a.logger.Error().Msg("LoadingDataCallback is nil")
		return
	}

	// load auth data...
	a.logger.Info().Msg("auth get from repo...")

	auth, err := a.authRepo.Get(ctx)

	if err != nil {
		a.logger.Warn().Err(err).Msg("auth get from repo failed")
	}

	if auth != nil {
		a.logger.Info().Interface("auth", auth).
			Msg("auth get from repo success, set to api client...")
		a.apiClient.SetAuthToken(auth.Token)
	} else {
		a.logger.Info().Msg("auth not set, show login screen...")
		a.showScreen(ScreenLogin)

		return
	}

	a.logger.Info().Msg("playerInfo loading starts...")
	callback.SendText("Загрузка информации о заведении...")

	playerInfo, err := a.playerInfoRepo.Get(ctx)

	if err != nil {
		a.logger.Warn().Err(err).Msg("playerInfo loading failed")

		switch err.(type) {
		default:
			callback.SendErrorMessage(err.Error())

		case *api.UnauthorizedError:
			a.showScreen(ScreenLogin)
		}

		return
	}

	a.logger.Info().Interface("playerInfo", playerInfo).Msg("playerInfo loaded")

	a.showScreen(ScreenPlayer)
}

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

func (a *App) iniLogger() zerolog.Logger {
	var output = a.config.LogWriter

	if a.config.Debug {
		output = zerolog.ConsoleWriter{Out: output}
	}

	logger := zerolog.New(output).
		With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)

	if a.config.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}

	return logger
}

func (a *App) showScreen(name string) {
	a.logger.Info().Str("name", name).Msg("show screen")

	if a.callbackMain == nil {
		a.logger.Error().Msg("MainCallback is nil")
		return
	}

	a.callbackMain.ShowScreen(name)
}

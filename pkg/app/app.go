package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/database"
	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/domain/repositories"
)

type App struct {
	config       Config
	callbackMain CallbackMain

	logger zerolog.Logger

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

func (a *App) Init() error {
	// connect db...
	db, err := database.Connect(a.config.DBFilePath)

	if err != nil {
		return fmt.Errorf("db connect fail: %w", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			a.logger.Warn().Err(err).Msg("db close failed")
		}
	}()

	a.logger.Info().Msg("database connected")

	// run migrate db...
	if err := database.Migrate(db); err != nil {
		return fmt.Errorf("database migrate: %w", err)
	}

	a.logger.Info().Msg("database migration applied")

	// init common...
	var apiClient api.Client = api.NewHTTPClient(a.config.ApiBaseURL)

	// init repos...
	a.playerInfoRepo = repositories.NewPlayerInfoApiRepo(apiClient)
	a.loginRepo = repositories.NewLoginApiRepo(apiClient)
	a.authRepo = repositories.NewAuthSqlLiteRepo(db)

	// show first screen...
	a.showScreen(ScreenLoadingData)

	return nil
}

func (a *App) LoadingData(ctx context.Context, callback CallbackLoadingData) {
	if callback == nil {
		a.logger.Error().Msg("LoadingDataCallback is nil")
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

	fmt.Println(playerInfo)

	return
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

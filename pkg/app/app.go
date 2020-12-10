package app

import (
	"path"

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

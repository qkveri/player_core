package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"

	"github.com/qkveri/player_core/pkg/api"
	"github.com/qkveri/player_core/pkg/domain"
	"github.com/qkveri/player_core/pkg/domain/repositories"
	"github.com/qkveri/player_core/pkg/state"
)

type App struct {
	config       Config
	callbackMain CallbackMain

	state     *state.State
	logger    zerolog.Logger
	apiClient api.Client

	// repos...
	playerInfoRepo domain.PlayerInfoRepository
	loginRepo      domain.LoginRepository
	musicDataRepo  domain.MusicDataRepository
	authRepo       domain.AuthRepository
}

func NewApp(config Config, callbackMain CallbackMain) *App {
	return &App{
		config:       config,
		callbackMain: callbackMain,
	}
}

func (a *App) Init() {
	// init state...
	a.state = state.NewState()

	// init logger...
	a.logger = a.iniLogger()

	// init common...
	a.apiClient = api.NewHTTPClient(a.config.ApiBaseURL)

	// init repos...
	a.playerInfoRepo = repositories.NewPlayerInfoApiRepo(a.apiClient)
	a.loginRepo = repositories.NewLoginApiRepo(a.apiClient)
	a.musicDataRepo = repositories.NewMusicDataApiRepo(a.apiClient)
	a.authRepo = repositories.NewAuthFileRepo(path.Join(a.config.DataDir, "a.tk"), a.config.SecretKey)

	// show first screen...
	a.showScreen(ScreenLoadingData)
}

func (a *App) Run(ctx context.Context) {
	<-ctx.Done()
}

func (a *App) iniLogger() zerolog.Logger {
	var output io.Writer

	if a.config.Debug {
		output = os.Stdout
	} else if file, err := a.createLogFile(); err != nil {
		log.New(os.Stderr, "LOGGER_CREATE_FILE: ", log.LstdFlags).
			Printf("create log file failed: %v\n", err)
	} else {
		output = file
	}

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

func (a *App) createLogFile() (*os.File, error) {
	dir := path.Join(a.config.CacheDir, "logs")

	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return nil, fmt.Errorf("os.Mkdir: %w, dir: %s", err, dir)
			}
		} else {
			return nil, fmt.Errorf("os.Stat: %w, dir: %s", err, dir)
		}
	}

	filePath := path.Join(dir, fmt.Sprintf("%s.txt", time.Now().Format(time.RFC3339)))

	file, err := os.Create(filePath)

	if err != nil {
		return nil, fmt.Errorf("os.Create: %w, filePath: %s", err, filePath)
	}

	return file, nil
}

func (a *App) showScreen(name string) {
	a.logger.Debug().Str("name", name).Msg("show screen")

	if a.callbackMain == nil {
		a.logger.Error().Msg("MainCallback is nil")
		return
	}

	a.callbackMain.ShowScreen(name)
}

func (a *App) errMessageForClient(err error) string {
	if a.config.Debug {
		return err.Error()
	}

	switch err.(type) {
	default:
		return "Произошла ошибка, пожалуйста, обратитесь в службу поддержки"

	case *api.NoInternetError:
		return "Нет подключения к интернету"
	}
}

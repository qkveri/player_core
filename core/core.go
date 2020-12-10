package core

import (
	"context"
	"os"

	"github.com/qkveri/player_core/pkg/app"
)

type (
	CallbackMain        interface{ app.CallbackMain }
	CallbackLoadingData interface{ app.CallbackLoadingData }
	CallbackLogin       interface{ app.CallbackLogin }
)

var (
	ctx, ctxCancel = context.WithCancel(context.Background())

	a *app.App

	// callbacks
	cbLoadingData CallbackLoadingData
	cbLogin       CallbackLogin
)

func InitApp(
	debug bool,
	secretKey string,
	apiBaseURL string,
	dataDir string,
	cacheDir string,

	callbackMain CallbackMain,
) {
	config := app.Config{
		Debug:     debug,
		LogWriter: os.Stdout,

		SecretKey:  secretKey,
		ApiBaseURL: apiBaseURL,

		DataDir:  dataDir,
		CacheDir: cacheDir,
	}

	a = app.NewApp(config, callbackMain)
	a.Init()
}

func Shutdown() {
	ctxCancel()
}

func RegisterLoadingDataCallback(callback CallbackLoadingData) {
	cbLoadingData = callback
}

func LoadingData() {
	a.LoadingData(ctx, cbLoadingData)
}

func RegisterLoginCallback(callback CallbackLogin) {
	cbLogin = callback
}

func Login(code string) {
	a.Login(ctx, cbLogin, code)
}

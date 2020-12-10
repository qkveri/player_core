package core

import (
	"context"

	"github.com/qkveri/player_core/pkg/app"
)

type (
	CallbackMain     interface{ app.CallbackMain }
	CallbackLoadData interface{ app.CallbackLoadData }
	CallbackLogin    interface{ app.CallbackLogin }
)

var (
	ctx, ctxCancel = context.WithCancel(context.Background())

	a *app.App

	// callbacks
	cbLoadData CallbackLoadData
	cbLogin    CallbackLogin
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
		Debug: debug,

		SecretKey:  secretKey,
		ApiBaseURL: apiBaseURL,

		DataDir:  dataDir,
		CacheDir: cacheDir,
	}

	a = app.NewApp(config, callbackMain)
	a.Init()
}

func Run() {
	a.Run(ctx)
}

func Shutdown() {
	ctxCancel()
}

func RegisterLoadDataCallback(callback CallbackLoadData) {
	cbLoadData = callback
}

func LoadData() {
	a.LoadData(ctx, cbLoadData)
}

func RegisterLoginCallback(callback CallbackLogin) {
	cbLogin = callback
}

func Login(code string) {
	a.Login(ctx, cbLogin, code)
}

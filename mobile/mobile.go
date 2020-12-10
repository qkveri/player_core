package mobile

import (
	"context"

	"github.com/qkveri/player_core/pkg/app"
)

type CallbackMain interface {
	app.CallbackMain
}

type CallbackLoadingData interface {
	app.CallbackLoadingData
}

var (
	ctx = context.Background()

	a *app.App

	// callbacks
	cbLoadingData CallbackLoadingData
)

func InitApp(
	debug bool,
	apiBaseURL string,
	authFilePath string,
	authKey string,

	callbackMain CallbackMain,
) {
	config := app.Config{
		Debug:        debug,
		ApiBaseURL:   apiBaseURL,
		AuthFilePath: authFilePath,
		AuthKey:      authKey,
	}

	a = app.NewApp(config, callbackMain)
}

func RegisterLoadingDataCallback(callback CallbackLoadingData) {
	cbLoadingData = callback
}

func LoadingData() {
	a.LoadingData(ctx, cbLoadingData)
}

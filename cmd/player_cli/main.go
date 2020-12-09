package main

// *****************
// * GoMobile like *
// *****************

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/qkveri/player_core/pkg/app"
)

var (
	ctx context.Context

	a *app.App

	// callbacks
	cbLoadingData app.CallbackLoadingData
	cbLogin       app.CallbackLogin
)

func init() {
	var ctxCancel context.CancelFunc

	ctx, ctxCancel = context.WithCancel(context.Background())

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop

		ctxCancel()
		signal.Stop(stop)
		close(stop)
	}()
}

func main() {
	// Вызывается при инициализации...
	cb := &callbackMain{}
	InitApp(cb)
}

// Global Init APP...
func InitApp(mainCallback app.CallbackMain) {
	config := app.Config{
		Debug: true,

		ApiBaseURL: "https://api.muzplat.ru/api/player/v2",
		DBFilePath: "/Users/petr/dev/apps/pult/player_core/app.db",
		LogWriter:  os.Stdout,
		//LogWriter:  new(bytes.Buffer),
	}

	a = app.NewApp(config, mainCallback)

	if err := a.Init(); err != nil {
		log.Fatalf("init app failed: %v", err)
	}
}

// For Loading Data Screen...
func RegisterLoadingDataCallback(callback app.CallbackLoadingData) {
	cbLoadingData = callback
}

func LoadingData() {
	a.LoadingData(ctx, cbLoadingData)
}

// For Login Screen...
func RegisterLoginCallback(callback app.CallbackLogin) {
	cbLogin = callback
}

func Login(code string) {
	a.Login(ctx, cbLogin, code)
}

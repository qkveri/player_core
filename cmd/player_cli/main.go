package main

// *****************
// * GoMobile like *
// *****************

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/qkveri/player_core/core"
)

func init() {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop

		core.Shutdown()
		signal.Stop(stop)
		close(stop)
	}()
}

func main() {
	cb := &callbackMain{}

	core.InitApp(
		false,
		"85dea59886138936d3b1a573f6069357",
		"https://api.muzplat.ru/api/player/v2",
		"/Users/petr/dev/apps/pult/player_core/tmp/data",
		"/Users/petr/dev/apps/pult/player_core/tmp/cache",
		cb,
	)

	core.Run()
}

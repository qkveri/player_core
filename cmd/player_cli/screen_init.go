package main

import (
	"log"

	"github.com/qkveri/player_core/pkg/app"
)

type callbackMain struct {
}

func (m *callbackMain) ShowScreen(name string) {
	switch name {
	default:
		log.Fatalf("unknown screen: %s", name)

	case app.ScreenLoadingData:
		openLoadingScreen()

	case app.ScreenLogin:
		openLoginScreen()
	}
}

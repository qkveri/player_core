package main

import (
	"fmt"

	"github.com/qkveri/player_core/core"
)

func openLoadingScreen() {
	cb := &callbackLoadingData{}

	core.RegisterLoadingDataCallback(cb)
	core.LoadingData()
}

type callbackLoadingData struct {
}

func (l *callbackLoadingData) SendText(text string) {
	fmt.Printf("💾 LoadingText: %s\n", text)
}

func (l *callbackLoadingData) SendErrorMessage(message string) {
	fmt.Printf("\n❌ Ошибка загрузки: %s\n", message)

	if waitConfirm("Повторить? [Y/n]: ") {
		core.LoadingData()
	}
}

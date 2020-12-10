package main

import (
	"fmt"

	"github.com/qkveri/player_core/core"
)

func openLoadingScreen() {
	cb := &callbackLoadData{}

	core.RegisterLoadDataCallback(cb)
	core.LoadData()
}

type callbackLoadData struct {
}

func (l *callbackLoadData) SendText(text string) {
	fmt.Printf("💾 LoadingText: %s\n", text)
}

func (l *callbackLoadData) SendErrorMessage(message string) {
	fmt.Printf("\n❌ Ошибка загрузки: %s\nПовторить? [Y/n]: ", message)

	if waitConfirm() {
		core.LoadData()
	}
}

func (l *callbackLoadData) SendPlayerInfo(json string) {
	fmt.Printf("💾 SetPlayerInfo: %s\n", json)
}

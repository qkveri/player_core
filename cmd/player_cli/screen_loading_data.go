package main

import (
	"fmt"
)

func openLoadingScreen() {
	cb := &callbackLoadingData{}
	RegisterLoadingDataCallback(cb)
	LoadingData()
}

type callbackLoadingData struct {
}

func (l *callbackLoadingData) SendText(text string) {
	fmt.Printf("💾 LoadingText: %s\n", text)
}

func (l *callbackLoadingData) SendErrorMessage(message string) {
	fmt.Printf("\n❌ Ошибка загрузки: %s\n", message)

	if waitConfirm("Повторить? [Y/n]: ") {
		LoadingData()
	}
}

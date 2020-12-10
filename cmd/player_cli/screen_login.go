package main

import (
	"fmt"

	"github.com/qkveri/player_core/core"
)

func openLoginScreen() {
	cb := &callbackLogin{}

	core.RegisterLoginCallback(cb)

	// нажали "Вход"
	fmt.Printf("\n🔒 Не авторизованы.\n")
	fmt.Print("Введите код: ")

	cb.code = waitInput()
	cb.login()
}

type callbackLogin struct {
	code string
}

func (l *callbackLogin) SendErrorMessage(message string) {
	fmt.Printf("\n❌ Ошибка логина: %s\nПовторить? [Y/n]:", message)

	if waitConfirm() {
		l.login()
	}
}

func (l *callbackLogin) SendCodeIncorrectErrorMessage(message string) {
	fmt.Printf("\n❌ Неверный код: %s\n", message)
	fmt.Print("Введите новый код: ")

	l.code = waitInput()
	l.login()
}

func (l *callbackLogin) login() {
	core.Login(l.code)
}

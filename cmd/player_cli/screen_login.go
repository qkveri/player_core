package main

import "fmt"

func openLoginScreen() {
	cb := &callbackLogin{}
	RegisterLoginCallback(cb)

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
	fmt.Printf("\n❌ Ошибка логина: %s\n", message)

	if waitConfirm("Повторить? [Y/n]: ") {
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
	Login(l.code)
}

package app

import "io"

const (
	ScreenLoadingData = "loading"
	ScreenLogin       = "login"
	ScreenPlayer      = "player"
)

type Config struct {
	Debug     bool
	LogWriter io.Writer

	SecretKey  string
	ApiBaseURL string

	DataDir  string
	CacheDir string
}

type CallbackMain interface {
	ShowScreen(name string)
}

type CallbackLoadData interface {
	SendText(text string)
	SendErrorMessage(message string)
	SendPlayerInfo(json string)
}

type CallbackLogin interface {
	SendErrorMessage(message string)
	SendCodeIncorrectErrorMessage(message string)
}

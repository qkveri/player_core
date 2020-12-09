package app

import (
	"io"
)

const (
	ScreenLoadingData = "loading"
	ScreenLogin       = "login"
	ScreenPlayer      = "player"
)

type Config struct {
	Debug bool

	ApiBaseURL string
	DBFilePath string
	LogWriter  io.Writer
}

// Main...
type CallbackMain interface {
	ShowScreen(name string)
}

// Loading data...
type CallbackLoadingData interface {
	SendText(text string)
	SendErrorMessage(message string)
}

// Login...
type CallbackLogin interface {
	SendErrorMessage(message string)
	SendCodeIncorrectErrorMessage(message string)
}

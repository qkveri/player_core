package app

const (
	ScreenLoadingData = "loading"
	ScreenLogin       = "login"
	ScreenPlayer      = "player"
)

type Config struct {
	Debug bool

	SecretKey  string
	ApiBaseURL string

	DataDir  string
	CacheDir string
}

type CallbackMain interface {
	ShowScreen(name string)
	SendErrorMessage(message string)
}

type CallbackLoadData interface {
	SendText(text string)
	SendErrorMessage(message string)
}

type CallbackLogin interface {
	SendErrorMessage(message string)
	SendCodeIncorrectErrorMessage(message string)
}

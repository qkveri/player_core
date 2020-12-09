package mobile

import (
	"fmt"
	"time"
)

type JavaCallback interface {
	SendString(string)
}

// callback
var jc JavaCallback

func RegisterJavaCallback(c JavaCallback) {
	jc = c
}

func Init(cachePath string) {
	fmt.Println(cachePath)
}

func Callback() {
	for {
		jc.SendString("tick")
		<-time.After(time.Second)
	}
}

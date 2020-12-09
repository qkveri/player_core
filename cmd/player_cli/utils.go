package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func waitInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return strings.TrimSpace(text)
}

func waitConfirm(message string) bool {
	fmt.Print(message)

	text := waitInput()
	text = strings.ToLower(text)

	switch text {
	case "n", "no":
		return false
	default:
		return true
	}
}

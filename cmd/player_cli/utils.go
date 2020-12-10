package main

import (
	"bufio"
	"os"
	"strings"
)

func waitInput() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	text = strings.ReplaceAll(text, "\n", "")

	return strings.TrimSpace(text)
}

func waitConfirm() bool {
	text := waitInput()
	text = strings.ToLower(text)

	switch text {
	case "n", "no":
		return false
	default:
		return true
	}
}

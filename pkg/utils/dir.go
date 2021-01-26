package utils

import (
	"fmt"
	"os"
)

func MkDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				return fmt.Errorf("os.Mkdir: %w", err)
			}
		} else {
			return fmt.Errorf("os.Stat: %w", err)
		}
	}

	return nil
}

package config

import (
	"fmt"
	"os"
)

// Read a file from a disk
func read(file string) []byte {
	content, err := os.ReadFile("./ops/configs/" + file)

	if err != nil {
		panic(fmt.Sprintf("error: %v", err))
	}

	return content
}

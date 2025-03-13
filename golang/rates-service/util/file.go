package util

import (
	"fmt"
	"os"
)

func WriteFile(dir, file string, data []byte) error {
	fullPath := fmt.Sprintf("%s/%s", dir, file)
	return os.WriteFile(fullPath, data, 0644)
}

package utils

import (
	"os"
	"path/filepath"
	"strconv"
)

func GetUniqueFileName(dir, baseName, ext string) string {
	filePath := filepath.Join(dir, baseName+ext)
	counter := 1
	for {
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return filePath
		}
		filePath = filepath.Join(dir, baseName+strconv.Itoa(counter)+ext)
		counter++
	}
}

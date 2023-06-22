package util

import (
	"path/filepath"
)

func GetFileExtension(filename string) string {
	extension := filepath.Ext(filename)
	//remove the dot
	extension = extension[1:]
	return extension
}

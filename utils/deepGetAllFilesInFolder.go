package utils

import (
	"os"
	"path/filepath"
)

func DeepGetAllFilesInFolder(absoluteFolderPath string) []string {
	var files []string
	err := filepath.Walk(absoluteFolderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	HandleErrorStrictly(err)
	return files
}

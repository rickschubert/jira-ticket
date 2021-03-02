package utils

import (
	"os"
	"path/filepath"
)

func getMockFeatureFileList() []string {
	return []string{"~/a/file/one.feature", "~/a/file/two.feature", "~/another/file/another.js", "~/organisations/organisationsViewer.feature"}
}

func DeepGetAllFilesInFolder(absoluteFolderPath string) []string {
	if os.Getenv("UNIT_TESTS") == "true" {
		return getMockFeatureFileList()
	}
	var files []string
	err := filepath.Walk(absoluteFolderPath, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	HandleErrorStrictly(err)
	return files
}

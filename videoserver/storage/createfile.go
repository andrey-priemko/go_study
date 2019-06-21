package storage

import (
	"os"
	"path/filepath"
)

func CreateFile(fileName string, id string) (*os.File, error) {
	if err := os.Mkdir(ContentDir, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	dirPath := filepath.Join(ContentDir, id)
	if err := os.Mkdir(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, err
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, err
}
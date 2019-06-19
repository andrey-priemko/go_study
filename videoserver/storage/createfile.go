package storage

import (
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

func CreateFile(fileName string) (*os.File, string, error) {
	if err := os.Mkdir(ContentDir, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, "", err
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, "", err
	}
	idStr := id.String()

	dirPath := filepath.Join(ContentDir, idStr)
	if err := os.Mkdir(dirPath, os.ModeDir); err != nil && !os.IsExist(err) {
		return nil, "", err
	}

	filePath := filepath.Join(dirPath, fileName)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	return file, idStr, err
}
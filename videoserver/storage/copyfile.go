package storage

import (
	"io"
	"mime/multipart"
)

func CopyFile(fileReader multipart.File, id string) error {
	file, err := CreateFile(VideoFileName, id)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, fileReader)
	if err != nil {
		return err
	}

	return nil
}

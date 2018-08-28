package utils

import (
	"fmt"
	"path/filepath"
	"path"
	"os"
	"github.com/satori/go.uuid"
	"../logger"
)


func GenPath(mediaPath, filename string) (string, string, uuid.UUID) {
	ext := path.Ext(filename)

	guid := uuid.Must(uuid.NewV4())
	sguid := fmt.Sprintf("%s", guid)

	newFileName := fmt.Sprintf("%s%s", guid, ext)
	relPath := filepath.Join(sguid[:2], sguid[2:4])
	absPath := filepath.Join(mediaPath, relPath)

	err := os.MkdirAll(absPath, 0777)
	if err != nil {
		logger.Error.Println(err)
	}

	pathToFile := filepath.Join(absPath, newFileName)
	relPathToFile := filepath.Join(relPath, newFileName)

	logger.Info.Printf("new file save: %s", pathToFile)
	return pathToFile, relPathToFile, guid
}
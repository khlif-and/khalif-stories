package utils

import (
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

)

func ConvertToAAC(inputFile multipart.File, originalFilename string) (*os.File, string, error) {
	tempInput, err := os.CreateTemp("", "input-*"+filepath.Ext(originalFilename))
	if err != nil {
		return nil, "", err
	}
	defer os.Remove(tempInput.Name())

	if _, err := io.Copy(tempInput, inputFile); err != nil {
		return nil, "", err
	}
	tempInput.Close()

	tempOutputName := tempInput.Name() + ".m4a"

	cmd := exec.Command("ffmpeg", "-i", tempInput.Name(), "-c:a", "aac", "-b:a", "128k", "-vn", "-y", tempOutputName)
	
	if err := cmd.Run(); err != nil {
		return nil, "", err
	}

	outputFile, err := os.Open(tempOutputName)
	if err != nil {
		return nil, "", err
	}

	return outputFile, tempOutputName, nil
}
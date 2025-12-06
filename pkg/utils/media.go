package utils

import (
	"bytes"
	"context"
	"mime/multipart"
	"path/filepath"

)

// UploadAndAnalyzeImage: Upload Gambar + Ekstrak Warna
func UploadAndAnalyzeImage(ctx context.Context, uploader *AzureUploader, file multipart.File, header *multipart.FileHeader, containerName, folderPath, fileUUID string) (string, string, error) {
	if file == nil {
		return "", "", nil
	}

	fileBytes, err := ReadMultipartFileToBytes(file)
	if err != nil {
		return "", "", err
	}

	if fileBytes == nil {
		return "", "", nil
	}

	filename := folderPath + fileUUID + filepath.Ext(header.Filename)
	
	imageURL, err := uploader.UploadToContainer(ctx, file, containerName, filename)
	if err != nil {
		return "", "", err
	}

	dominantColor := "#000000"
	if color, err := ExtractDominantColor(bytes.NewReader(fileBytes)); err == nil {
		dominantColor = color
	}

	return imageURL, dominantColor, nil
}

// BARU: UploadFile (Generic untuk Audio/File lain tanpa analisis warna)
func UploadFile(ctx context.Context, uploader *AzureUploader, file multipart.File, header *multipart.FileHeader, containerName, folderPath, fileUUID string) (string, error) {
	if file == nil {
		return "", nil
	}

	// Gunakan helper baca file yg sama
	fileBytes, err := ReadMultipartFileToBytes(file)
	if err != nil {
		return "", err
	}
	
	if fileBytes == nil {
		return "", nil
	}

	filename := folderPath + fileUUID + filepath.Ext(header.Filename)
	
	return uploader.UploadToContainer(ctx, file, containerName, filename)
}
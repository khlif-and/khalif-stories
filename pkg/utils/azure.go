package utils

import (
	"context"
	"mime/multipart"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

)

type AzureUploader struct {
	Client        *azblob.Client
	ContainerName string // Ini container default (profile-pick)
}

func NewAzureUploader(connStr, containerName string) (*AzureUploader, error) {
	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	return &AzureUploader{Client: client, ContainerName: containerName}, nil
}

// UploadFile menggunakan container default dari config (profile-pick)
func (a *AzureUploader) UploadFile(file multipart.File, filename string) (string, error) {
	return a.UploadToContainer(file, a.ContainerName, filename)
}

// UploadToContainer memungkinkan kita memilih container tujuan (misal: "category")
func (a *AzureUploader) UploadToContainer(file multipart.File, containerName, filename string) (string, error) {
	ctx := context.Background()
	// UploadStream ke container spesifik
	_, err := a.Client.UploadStream(ctx, containerName, filename, file, nil)
	if err != nil {
		return "", err
	}
	return "https://" + a.Client.URL() + "/" + containerName + "/" + filename, nil
}
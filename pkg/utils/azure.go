package utils

import (
	"context"
	"mime/multipart"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

)

type AzureUploader struct {
	Client        *azblob.Client
	ContainerName string
}

func NewAzureUploader(connStr, containerName string) (*AzureUploader, error) {
	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		return nil, err
	}
	return &AzureUploader{Client: client, ContainerName: containerName}, nil
}

func (a *AzureUploader) UploadFile(file multipart.File, filename string) (string, error) {
	return a.UploadToContainer(file, a.ContainerName, filename)
}

func (a *AzureUploader) UploadToContainer(file multipart.File, containerName, filename string) (string, error) {
	ctx := context.Background()
	_, err := a.Client.UploadStream(ctx, containerName, filename, file, nil)
	if err != nil {
		return "", err
	}

	baseURL := a.Client.URL()
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL + containerName + "/" + filename, nil
}
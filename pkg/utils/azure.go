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

func (a *AzureUploader) Upload(file multipart.File, header *multipart.FileHeader) (string, error) {
	return a.UploadToContainer(context.Background(), file, a.ContainerName, header.Filename)
}

func (a *AzureUploader) Delete(fileURL string) error {
	parts := strings.Split(fileURL, "/")
	if len(parts) == 0 {
		return nil
	}
	blobName := parts[len(parts)-1]
	_, err := a.Client.DeleteBlob(context.Background(), a.ContainerName, blobName, nil)
	return err
}

func (a *AzureUploader) UploadFile(ctx context.Context, file multipart.File, filename string) (string, error) {
	return a.UploadToContainer(ctx, file, a.ContainerName, filename)
}

func (a *AzureUploader) UploadToContainer(ctx context.Context, file multipart.File, containerName, filename string) (string, error) {
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
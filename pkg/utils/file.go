package utils

import (
	"bytes"
	"io"
	"mime/multipart"

)

// ReadMultipartFileToBytes membaca file ke buffer dan mereset pointer agar siap di-upload
func ReadMultipartFileToBytes(file multipart.File) ([]byte, error) {
	if file == nil {
		return nil, nil
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		return nil, err
	}

	// Reset file pointer ke awal agar bisa dibaca ulang oleh uploader
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
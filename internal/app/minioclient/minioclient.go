package minioclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient представляет клиент MinIO.
type MinioClient struct {
	Client *minio.Client
}

// NewMinioClient создает новый экземпляр клиента MinIO.
func NewMinioClient() (*MinioClient, error) {
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	endpoint  := os.Getenv("MINIO_ENDPOINT")
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}

	return &MinioClient{
		Client: client,
	}, nil
}

// UploadServiceImage загружает изображение в MinIO и возвращает URL изображения.
func (mc *MinioClient) UploadServiceImage(baggageID int, imageBytes []byte, contentType string) (string, error) {
	objectName := fmt.Sprintf("baggages/%d/image", baggageID)

	reader := io.NopCloser(bytes.NewReader(imageBytes))

	_, err := mc.Client.PutObject(context.TODO(), "images-bucket", objectName, reader, int64(len(imageBytes)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", errors.New("ошибка загрузки изображения в минио")
	}

	// Формирование URL изображения
	imageURL := fmt.Sprintf("http://localhost:9000/images-bucket/%s", objectName)
	return imageURL, nil
}

// RemoveServiceImage удаляет изображение услуги из MinIO.
func (mc *MinioClient) RemoveServiceImage(baggageID int) error {
	objectName := fmt.Sprintf("baggages/%d/image", baggageID)
	err := mc.Client.RemoveObject(context.TODO(), "images-bucket", objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return errors.New("не удалось удалить изображение из бакет")
	}
	return nil
}

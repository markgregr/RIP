package main

import (
	"context"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

//mc config host add myminio http://localhost:9000 yD8s9f2R1tU3qA7mP X6bN1oY9zP4wG2cV5eF8hR2tA7mS3kL1
//mc mb myminio/ImagesBucket
//http://yD8s9f2R1tU3qA7mP:X6bN1oY9zP4wG2cV5eF8hR2tA7mS3kL1@localhost:9000/images-bucket/1.png.
func main() {
	ctx := context.Background()
	endpoint := "localhost:9000"
	accessKeyID := "yD8s9f2R1tU3qA7mP"
	secretAccessKey := "X6bN1oY9zP4wG2cV5eF8hR2tA7mS3kL1"
	useSSL := false

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	imageFolder := "./resources/images"
	imageFileNames := []string{"1.png", "2.png", "3.png", "4.png", "5.png", "6.png"}
	bucketName := "images-bucket"

	for _, fileName := range imageFileNames {
		file, err := os.Open(imageFolder + "/" + fileName)
		if err != nil {
			log.Println("Error opening file:", fileName, err)
			continue
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Println("Error getting file info:", fileName, err)
			continue
		}

		objectName := fileName
		n, err := client.PutObject(ctx, bucketName, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
			ContentType: "image/png",
		})
		if err != nil {
			log.Println("Error uploading file:", fileName, err)
			continue
		}

		log.Println("Successfully uploaded", objectName, "of size:", n, "bytes.")
	}
}

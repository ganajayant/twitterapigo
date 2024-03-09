package initializers

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var Bucket *storage.BucketHandle

func ConnectFireBase() {
	bucketurl := os.Getenv("Bucket")
	config := &firebase.Config{
		StorageBucket: bucketurl,
	}
	FilePath := os.Getenv("FilePath")
	opt := option.WithCredentialsFile(FilePath)
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error connecting Firebase")
	}

	client, err := app.Storage(context.Background())
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error connecting Firebase")
	}
	Bucket, err = client.DefaultBucket()
	if err != nil {
		log.Fatalln("Error connecting Firebase Bucket")
	}
}

func UploadFile(file *multipart.FileHeader) (string, error) {
	imagepath := file.Filename
	wc := Bucket.Object(imagepath).NewWriter(context.Background())
	f, err := file.Open()
	if err != nil {
		return "", err
	}

	_, err = io.Copy(wc, f)
	if err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}
	imageurl := "https://storage.cloud.google.com/" + os.Getenv("Bucket") + "/" + imagepath
	return imageurl, nil
}

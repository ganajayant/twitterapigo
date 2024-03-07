package initializers

import (
	"context"
	"fmt"
	"log"
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
	opt := option.WithCredentialsFile("gotwitter-38aeb-firebase-adminsdk-blu76-0fa61a6c13.json")
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

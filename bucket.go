package main

import (
	"context"
	"log"
	"os"
	"time"

	gstorage "cloud.google.com/go/storage"
	goauth2 "golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	goption "google.golang.org/api/option"
)

func listBucketFiles() {
	credentialsFilePath := os.Getenv("CREDENTIALS_FILE_PATH")
	if credentialsFilePath == "" {
		log.Fatal("no CREDENTIALS_FILE_PATH env var found")
	}
	credentialsFileData, err := os.ReadFile(credentialsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	credentials, err := goauth2.CredentialsFromJSON(context.TODO(), credentialsFileData, gstorage.ScopeFullControl)
	if err != nil {
		log.Fatal(err)
	}
	opts := []goption.ClientOption{goption.WithCredentials(credentials)}

	gcsClient, err := gstorage.NewClient(context.TODO(), opts...)
	if err != nil {
		log.Fatal(err)
	}

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatal("no BUCKET_NAME env var found")
	}
	bucket := gcsClient.Bucket(bucketName)

	query := &gstorage.Query{Prefix: ""}
	var names []string
	it := bucket.Objects(context.TODO(), query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		names = append(names, attrs.Name)
	}
	log.Print(names)
}

func main() {
	for {
		listBucketFiles()
		time.Sleep(15 * time.Second)
	}
}

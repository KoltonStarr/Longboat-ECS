package main

import (
	// What does this native package do?
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

// TODO: Figure out if the upload logic works with other image formats.
// TODO: Implement a GET route that will get a single image from the bucket only if the request has the exact filename.
// TODO: Implements tests!
// TODO: Figure out a better way to provide the s3 client with configuration credentials so that this server can run in a container and have all it needs.
//       I might need to configure an IAM role for the application itself and give it access to only S3.
// TODO: Add proper error handling and response handling.

var MAX_BYTE_SIZE = int64(10 << 20)

func main() {
	router := chi.NewRouter()
	router.Use(LogPretty)
	router.Post("/images", imageUpload)

	fmt.Println("Server is running on port 3000")
	http.ListenAndServe(":3000", router)
}

func imageUpload(responseWriter http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(MAX_BYTE_SIZE)

	// Wuh oh file too big!
	if err != nil {
		color.RedString(err.Error())
	}

	// What IS "file" here?
	file, partHeader, _ := request.FormFile("file")
	fileName := partHeader.Filename
	defer file.Close()

	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	client := s3.NewFromConfig(cfg)

	_, s3Err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("hwot-bayrn"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if s3Err != nil {
		fmt.Println(err.Error())
	}
}

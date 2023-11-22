package main

/*
Note:

From what I have read in the AWS docs I might not have to do anything at all in order for my client S3 instantiation to work.
The LoadDefaulConfig function is apparently smart enough to use the IAM role of the ECS task that is running my container.
The only thing I need to confirm is that the IAM role of the task has the right permissions to actually write to my bucket.
*/

import (
	// What does this native package do?
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

// TODO: Implement a GET route that will get a single image from the bucket only if the request has the exact filename.
// TODO: Implements tests!
// TODO: Figure out how to do all of this with CloudFormation.
// TODO: Turn this into an HTTPS server and figure out how that works.
// TODO: Make it so that I can dynamically look for all file keys and parse each file without explicitly hard-coding file keys.

var MAX_BYTE_SIZE = int64(10 << 20)
var PORT = "8080"

func main() {
	router := chi.NewRouter()
	router.Use(LogPretty)
	router.Post("/images", imageUpload)

	fmt.Println("Server is running on port " + PORT)
	http.ListenAndServe(":"+PORT, router)
}

func imageUpload(w http.ResponseWriter, request *http.Request) {
	fileTooBigErr := request.ParseMultipartForm(MAX_BYTE_SIZE)
	if fileTooBigErr != nil {
		logAndWriteError(w, fileTooBigErr, http.StatusRequestEntityTooLarge)
		return
	}

	// What IS "file" here?
	imageFile, fileHeader, formFileErr := request.FormFile("file")
	if formFileErr != nil {
		logAndWriteError(w, formFileErr, http.StatusInternalServerError)
		return
	}

	fileName := fileHeader.Filename
	defer imageFile.Close()

	s3Client, s3ClientErr := createS3Client()
	if s3ClientErr != nil {
		logAndWriteError(w, s3ClientErr, http.StatusInternalServerError)
		return
	}

	_, s3OpErr := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String("hwot-bayrn"),
		Key:    aws.String(fileName),
		Body:   imageFile,
	})

	if s3OpErr != nil {
		logAndWriteError(w, s3OpErr, http.StatusInternalServerError)
	} else {
		log.Println(color.HiGreenString("** Images were uploaded successfully **"))
	}
}

func logAndWriteError(w http.ResponseWriter, err error, httpStatus int) {
	logError(err)
	w.WriteHeader(httpStatus)
	w.Write([]byte(err.Error()))
}

func createS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	client := s3.NewFromConfig(cfg)

	return client, err
}

func logError(err error) {
	log.Println(color.RedString(err.Error()))
}

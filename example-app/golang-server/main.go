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
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
)

// TODO: Finish the GET route.
// TODO: Implements tests!
// TODO: Figure out how to do all of this with CloudFormation.
// TODO: Turn this into an HTTPS server and figure out how that works.
// TODO: Make it so that I can dynamically look for all file keys and parse each file without explicitly hard-coding file keys.
// TODO: Write some logic that prints information about the image to the console.

var maxByteSize = int64(10 << 20)
var port = "8080"

func main() {
	r := chi.NewRouter()
	r.Use(LogPretty)

	r.Route("/images", func(r chi.Router) {
		r.Post("/", imageUploadHandler)
		r.Get("/{fileName}", getImageHandler)
	})

	fmt.Println("Server is running on port " + port)
	http.ListenAndServe(":"+port, r)
}

func getImageHandler(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")

	s3Client, s3ClientErr := createS3Client()
	if s3ClientErr != nil {
		logAndWriteError(w, s3ClientErr, http.StatusInternalServerError)
		return
	}

	output, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("hwot-bayrn"),
		Key:    aws.String(fileName),
	})

	if err != nil {
		logAndWriteError(w, err, http.StatusInternalServerError)
		return
	}
	// body here is a ReadCloser which essentially encapsulates a data source. In this case... the underlying data is going to be raw
	// binary data that represents the image I just got from S3. It isn't as simple as JSON. It is actual binary data.
	// To do something with it I need to read it into a slice of bytes.
	body := output.Body
	defer body.Close()
	imageData, _ := io.ReadAll(body)

	w.Header().Set("Content-Type", "image/png")
	w.Write(imageData)

}

func imageUploadHandler(w http.ResponseWriter, request *http.Request) {
	fileTooBigErr := request.ParseMultipartForm(maxByteSize)
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

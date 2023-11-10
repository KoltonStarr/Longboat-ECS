package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"context"
	// "log"
	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/s3"
)

// png, jpg
// var images = []image.Image{}

const uploadDir = "./image-uploads"

func main() {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	fmt.Println(cfg.Credentials)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/images", getImages)
	router.Post("/images", imageUpload)

	fmt.Println("Server is running on localhost 3000!")
	http.ListenAndServe(":3000", router)
}

func getImages(writer http.ResponseWriter, reader *http.Request) {
	fmt.Println("Hello!!!")
	byteData := []byte("Hello World!")
	writer.Write(byteData)
}

func imageUpload(writer http.ResponseWriter, reader *http.Request) {
	reader.ParseMultipartForm(10 << 20)
	file, handler, _ := reader.FormFile("file")
	defer file.Close()

	fileName := filepath.Join(uploadDir, handler.Filename)
	newFile, _ := os.Create(fileName)
	defer newFile.Close()

	io.Copy(newFile, file)
}

package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"time"

	"github.com/google/wire"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type IGoogleService interface {
	UploadFile(ctx context.Context, input model.ProfileImage) (string, error)
}
type GoogleService struct {
}

// define provider
var NewGoogleService = wire.NewSet(wire.Struct(new(GoogleService), "*"), wire.Bind(new(IGoogleService), new(*GoogleService)))

func (g *GoogleService) UploadFile(ctx context.Context, input model.ProfileImage) (string, error) {
	filename := fmt.Sprintf("%v-%s", time.Now().Format("01-02-2006"), *input.FileName)
	srv, err := drive.NewService(ctx, option.WithCredentialsFile("credential.json"), option.WithScopes(drive.DriveScope))
	if err != nil {
		return "", err
	}
	stream, err := base64.StdEncoding.DecodeString(*input.FileBase64)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return "", err
	}

	fileErr := ioutil.WriteFile(filename, stream, 0644)
	if fileErr != nil {
		fmt.Printf("file err %v", fileErr)
	}
	file, err := os.Open(filename)
	info, _ := file.Stat()
	if err != nil {
		log.Fatalf("Warning: %v", err)
	}

	if err != nil {
		log.Fatalln(err)
	}
	// Create File metadata
	cfg, _ := ctx.Value(entities.ConfigKey).(entities.PbConfig)
	fmt.Println("parent file id", cfg.ClientEmail)
	f := &drive.File{
		Parents: []string{"1otztb9RgRJFu-ObStqPrRnWNuM8GAoGC"},
		Name:    info.Name()}

	// Create and upload the file
	res, err := srv.Files.
		Create(f).
		Media(file).
		ProgressUpdater(func(now, size int64) { fmt.Printf("%d, %d\r", now, size) }).
		Do()
	file.Close()
	e := os.Remove(filename)
	if e != nil {
		fmt.Println("Error delete file from google, ", e)
	}
	if err != nil {
		fmt.Println("Error response from google, ", err)
	}
	fmt.Println("response from google, ", res.Id)
	return res.Id, nil
}

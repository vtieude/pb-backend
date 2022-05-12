package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"pb-backend/entities"
	"pb-backend/graph/model"
	"time"

	"github.com/google/wire"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

type IGoogleService interface {
	UploadFile(ctx context.Context, input model.ProfileImage) (string, error)
}
type GoogleService struct {
}

// define provider
var NewGoogleService = wire.NewSet(wire.Struct(new(GoogleService), "*"), wire.Bind(new(IGoogleService), new(*GoogleService)))

func (g *GoogleService) UploadFile(ctx context.Context, input model.ProfileImage) (string, error) {
	baseMimeType := "*/*"           // MimeType
	client := g.serviceAccount(ctx) // Please set the json file of Service account.
	filename := fmt.Sprintf("%v-%v", time.Now(), input.File.Filename)
	srv, err := drive.New(client)
	if err != nil {
		return "", err
	}
	stream, err := ioutil.ReadAll(input.File.File)

	if err != nil {
		return "", err
	}

	fileErr := ioutil.WriteFile(filename, stream, 0644)
	if fileErr != nil {
		fmt.Printf("file err %v", fileErr)
	}
	file, openErr := os.Open(filename)
	if openErr != nil {
		fmt.Printf("Error opening file: %v", openErr)
	}

	fileInf, err := file.Stat()
	if err != nil {
		return "", err
	}
	defer file.Close()
	cfg, _ := ctx.Value(entities.ConfigKey).(entities.PbConfig)
	f := &drive.File{
		Parents: []string{cfg.FolderId},
		Name:    filename}
	res, err := srv.Files.
		Create(f).
		ResumableMedia(context.Background(), file, fileInf.Size(), baseMimeType).
		ProgressUpdater(func(now, size int64) { fmt.Printf("%d, %d\r", now, size) }).
		Do()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", res.Id)
	return res.Id, nil
}

// ServiceAccount : Use Service account
func (g *GoogleService) serviceAccount(ctx context.Context) *http.Client {
	cfg, _ := ctx.Value(entities.ConfigKey).(entities.PbConfig)
	config := &jwt.Config{
		Email:      cfg.ClientEmail,
		PrivateKey: []byte(cfg.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
		TokenURL: google.JWTTokenURL,
	}
	client := config.Client(oauth2.NoContext)
	return client
}

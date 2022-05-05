package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"pb-backend/entities"

	"github.com/google/wire"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

type IGoogleService interface {
	UploadImage(ctx context.Context, fileName string) (bool, error)
}
type GoogleService struct {
}

// define provider
var NewGoogleService = wire.NewSet(wire.Struct(new(GoogleService), "*"), wire.Bind(new(IGoogleService), new(*GoogleService)))

func (g *GoogleService) UploadImage(ctx context.Context, fileName string) (bool, error) {
	filename := "tet.png"           // Filename
	baseMimeType := "*/*"           // MimeType
	client := g.serviceAccount(ctx) // Please set the json file of Service account.

	srv, err := drive.New(client)
	if err != nil {
		return false, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	fileInf, err := file.Stat()
	if err != nil {
		return false, err
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
	return true, nil
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

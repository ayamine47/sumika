package cloud

import (
	"context"
	"log"
	"os"

	"github.com/ayamine47/sumika/lib/config"
	"github.com/studio-b12/gowebdav"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var srv *drive.Service
var client *gowebdav.Client

func InitGoogleDrive() {
	secret, err := os.ReadFile("secret.json")
	if err != nil {
		log.Fatal("Failed to read client secret: ", err)
	}

	cfg, err := google.ConfigFromJSON(secret, drive.DriveFileScope, drive.DriveScope)
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}

	client := getClient(cfg)

	_, err = drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Failed to create google service: ", err)
	}

	log.Print("Initialized GoogleDrive")
}

func InitNextCloud() {
	client = gowebdav.NewClient(config.CurrentConfig.NextCloud.Url, config.CurrentConfig.NextCloud.Username, config.CurrentConfig.NextCloud.Password)
	err := client.Connect()
	if err != nil {
		log.Fatal("Failed to connect to NextCloud: ", err)
	}

	log.Print("Initialized NextCloud")
}

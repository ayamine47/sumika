package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ayamine47/sumika/lib/embed"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func UploadFileToGoogleDrive(fileName string, msgInfo *embed.MsgInfo) {
	file, err := os.Open(fileName + "_720.mp4")
	if err != nil {
		log.Print("Failed to open file: ", err)
	}

	defer file.Close()

	fileMetaData := &drive.File{
		Name:     fileName + ".mp4",
		MimeType: "video/mp4",
	}

	driveFile, err := srv.Files.Create(fileMetaData).Media(file).Do()
	if err != nil {
		log.Print("Failed to create file: ", err)
	}

	_, err = srv.Permissions.Create(driveFile.Id, &drive.Permission{Type: "anyone", Role: "reader"}).Do()
	if err != nil {
		log.Print("Failed to change file permission: ", err)
	}

	shareLink := fmt.Sprintf("https://drive.google.com/file/d/%s/view?usp=sharing", driveFile.Id)

	embed.SendSuccessEmbed(msgInfo, "Video URL: "+shareLink, msgInfo.OrgMsg.Reference())

	err = msgInfo.Session.MessageReactionRemove(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, "🤔", msgInfo.Session.State.User.ID)
	if err != nil {
		log.Print("Failed to add reaction: ", err)
	}

	err = msgInfo.Session.MessageReactionAdd(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, "✅")
	if err != nil {
		log.Print("Failed to add reaction: ", err)
	}
}

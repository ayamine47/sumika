package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/go-ytdlp"
	"github.com/u2takey/go-utils/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"

	ffmpeg "github.com/u2takey/ffmpeg-go"

	"github.com/ayamine47/sumika/lib/config"
)

var srv *drive.Service

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
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

// Retrieves a token from a local file.
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

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ytdlp.MustInstall(context.TODO(), nil)

	secret, err := os.ReadFile("secret.json")
	if err != nil {
		log.Fatal("Failed to read client secret: ", err)
	}

	cfg, err := google.ConfigFromJSON(secret, drive.DriveFileScope, drive.DriveScope)
	if err != nil {
		log.Fatal("Failed to create config: ", err)
	}

	client := getClient(cfg)

	srv, err = drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatal("Failed to create google service: ", err)
	}

	var token = config.CurrentConfig.Discord.Token
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create Bot")
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) { go MessageCreate(s, m) })
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent)

	err = discord.Open()
	if err != nil {
		log.Fatal("Failed to open Discord: ", err)
	}

	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)

	log.Print("sumika Running!")

	defer discord.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Shutdown...")
}

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("Recovering fatal error: ", err)
			debug.PrintStack()
		}
	}()

	//Ignore all messages created by the every bot
	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" || orgMsg.Author.Bot || orgMsg.ChannelID != config.CurrentConfig.Discord.Channel {
		return
	}

	url, err := url.Parse(orgMsg.Content)
	if err != nil {
		log.Print("Failed to parse URL: ", orgMsg.Content)
	}

	if url.Host != "youtu.be" && url.Host != "youtube.com" && url.Host != "www.youtube.com" && url.Host != "www.nicovideo.jp" {
		return
	}

	err = session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🤔")
	if err != nil {
		log.Print("Failed to add reaction: ", err)
	}

	id := uuid.NewUUID()

	dl := ytdlp.New().FormatSort("res,ext:mp4:m4a").Output(id + ".%(ext)s")

	res, err := dl.Run(context.TODO(), url.String())
	if err != nil {
		log.Print("Failed to download youtube video: ", err)
	}

	_ = res.OutputLogs

	log.Print("Downloaded")

	cmd := ffmpeg.Input(id+".mp4").Output(id+"_720.mp4", ffmpeg.KwArgs{"vf": "scale=1280:720,fps=30"})
	err = cmd.Run()
	if err != nil {
		log.Print("Failed to re-encode video: ", err.Error())
		err = session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "❌")
		if err != nil {
			log.Print("Failed to add reaction: ", err)
		}

		err = session.MessageReactionRemove(orgMsg.ChannelID, orgMsg.ID, "🤔", session.State.User.ID)
		if err != nil {
			log.Print("Failed to add reaction: ", err)
		}
		return
	}

	log.Print("re-encoded")

	file, err := os.Open(id + "_720.mp4")
	if err != nil {
		log.Print("Failed to open file: ", err)
	}

	defer file.Close()

	fileMetaData := &drive.File{
		Name:     id + "_720.mp4",
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

	_, err = session.ChannelMessageSendReply(orgMsg.ChannelID, shareLink, orgMsg.Reference())
	if err != nil {
		log.Print("Failed to send message: ", err)
	}

	err = session.MessageReactionRemove(orgMsg.ChannelID, orgMsg.ID, "🤔", session.State.User.ID)
	if err != nil {
		log.Print("Failed to add reaction: ", err)
	}

	err = session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "✅")
	if err != nil {
		log.Print("Failed to add reaction: ", err)
	}
}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/lrstanley/go-ytdlp"

	"github.com/ayamine47/sumika/lib/cloud"
	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/handler"
)

func main() {
	ytdlp.MustInstall(context.TODO(), nil)

	if config.CurrentConfig.NextCloud.Enable {
		cloud.InitNextCloud()
	}

	if config.CurrentConfig.GoogleDrive.Enable {
		cloud.InitGoogleDrive()
	}

	var token = config.CurrentConfig.Discord.Token
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Failed to create Bot")
	}

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) { go handler.MessageCreate(s, m) })
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent)

	err = discord.Open()
	if err != nil {
		log.Fatal("Failed to open Discord: ", err)
	}

	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)

	log.Print("sumika is now Running!")

	defer discord.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Shutdown...")
}

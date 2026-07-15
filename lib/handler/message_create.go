package handler

import (
	"log"
	"runtime/debug"
	"strings"

	"github.com/ayamine47/sumika/lib/cmds"
	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/embed"
	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	defer func() {
		if err := recover(); err != nil {
			log.Print("Recovering fatal error: ", err)
			debug.PrintStack()
		}
	}()

	//Ignore all messages created by the every bot
	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" || orgMsg.Author.Bot || orgMsg.ChannelID != config.CurrentConfig.Guild.Channel {
		return
	}

	msgInfo := embed.MsgInfo{
		Session: session,
		OrgMsg:  orgMsg,
	}

	if strings.HasPrefix(orgMsg.Content, config.CurrentConfig.Guild.Prefix) {
		cmd := strings.SplitN(strings.SplitN(orgMsg.Content, "\n", 2)[0], " ", 2)[0][len(config.CurrentConfig.Guild.Prefix):]
		switch cmd {
		case cmds.Ping:
			cmds.PingCmd(&msgInfo)
		case cmds.Help:
			cmds.HelpCmd(&msgInfo)
		case cmds.Get:
			cmds.GetCmd(&msgInfo)
		default:
			embed.SendErrorEmbed(&msgInfo, "No such command")
			return
		}
	}
}

package cmds

import (
	"runtime"
	"strconv"

	"github.com/ayamine47/sumika/lib/embed"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Ping = "ping"

func PingCmd(msgInfo *embed.MsgInfo) {
	msg := embed.NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(Ping)
	msg.Description = "Pong!"
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Golang",
		Value: "`" + runtime.GOARCH + " " + runtime.GOOS + " " + runtime.Version() + "`",
	})
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Stats",
		Value: "```\n" + strconv.Itoa(runtime.NumCPU()) + " cpu(s),\n" + strconv.Itoa(runtime.NumGoroutine()) + " go routine(s).```",
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Memory",
		Value: "```\n" + strconv.FormatUint(mem.Sys/1024/1024, 10) + "MB used,\n" + strconv.FormatUint(uint64(mem.NumGC), 10) + " GCs.```",
	})
	msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
}

package cmds

import (
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/embed"
	ytdlp "github.com/ayamine47/sumika/lib/yt-dlp"
)

const Get = "get"

func GetCmd(msgInfo *embed.MsgInfo) {
	split := strings.Split(msgInfo.OrgMsg.Content, " ")
	if len(split) > 1 {
		url, err := url.Parse(split[1])
		if err != nil {
			embed.SendErrorEmbed(msgInfo, "Failed to parse URL: "+"`"+split[1]+"`"+"\n"+err.Error())
			return
		}

		if !slices.Contains(config.CurrentConfig.UrlWhiteList, url.Host) {
			return
		}

		err = msgInfo.Session.MessageReactionAdd(msgInfo.OrgMsg.ChannelID, msgInfo.OrgMsg.ID, "🤔")
		if err != nil {
			log.Print("Failed to add reaction: ", err)
		}

		ytdlp.GetVideo(url, msgInfo)
	}
}

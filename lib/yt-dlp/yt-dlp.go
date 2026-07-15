package ytdlp

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/ayamine47/sumika/lib/embed"
	"github.com/ayamine47/sumika/lib/ffmpeg"
	"github.com/ayamine47/sumika/lib/utils"
	"github.com/lrstanley/go-ytdlp"
)

func GetVideo(u *url.URL, msgInfo *embed.MsgInfo) {
	metaDl := ytdlp.New().SkipDownload().DumpJSON()
	metaRes, err := metaDl.Run(context.TODO(), u.String())
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "Failed fetch metadata: "+"`"+u.String()+"`"+"\n"+err.Error())
		return
	}

	var meta utils.VideoMetadata

	err = json.Unmarshal([]byte(metaRes.Stdout), &meta)
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "Failed parse metadata: "+"`"+u.String()+"`"+"\n"+err.Error())
		return
	}

	dl := ytdlp.New().FormatSort("res,ext:mp4:m4a").Output("%(id)s_%(title)s.%(ext)s")

	_, err = dl.Run(context.TODO(), u.String())
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "Failed to download video: "+"`"+u.String()+"`"+"\n"+err.Error())
		return
	}

	fileName := meta.ID + "_" + meta.Title

	ffmpeg.ReEncodeVideo(fileName, msgInfo)
}

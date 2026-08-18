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
	cookie := getCookiePath(u)

	metaDl := ytdlp.New().SkipDownload().DumpJSON().JsRuntimes("nodejs")

	if cookie != "" {
		metaDl = metaDl.Cookies(cookie)
	}

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

	fileName := utils.SanitizeFileName(meta.ID + "_" + meta.Title)

	dl := ytdlp.New().FormatSort("res,ext:mp4:m4a").Output(fileName + ".%(ext)s").JsRuntimes("nodejs")
	if cookie != "" {
		dl = dl.Cookies(cookie)
	}

	_, err = dl.Run(context.TODO(), u.String())
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "Failed to download video: "+"`"+u.String()+"`"+"\n"+err.Error())
		return
	}

	ffmpeg.ReEncodeVideo(fileName, msgInfo)
}

package ffmpeg

import (
	"github.com/ayamine47/sumika/lib/cloud"
	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/embed"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReEncodeVideo(fileName string, msgInfo *embed.MsgInfo) {
	cmd := ffmpeg.Input(fileName+".mp4").Output(fileName+"_720.mp4", ffmpeg.KwArgs{"c:v": "h264", "c:a": "aac", "ab": "192k", "s": "1280x720", "r": "30", "vf": "scale=w=iw:h=ih"})
	err := cmd.Run()
	if err != nil {
		embed.SendErrorEmbed(msgInfo, "Failed to re-encode video: "+"`"+fileName+"`"+"\n"+err.Error())
		return
	}

	if config.CurrentConfig.GoogleDrive.Enable {
		cloud.UploadFileToGoogleDrive(fileName, msgInfo)
	}

	if config.CurrentConfig.NextCloud.Enable {
		cloud.UploadFileToNextCloud(fileName, msgInfo)
	}
}

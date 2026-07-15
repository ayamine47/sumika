package cmds

import (
	"github.com/ayamine47/sumika/lib/config"
	"github.com/ayamine47/sumika/lib/embed"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Help = "help"

func HelpCmd(msgInfo *embed.MsgInfo) {
	msg := embed.NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(Help)
	msg.Description = "# 使い方\n指定されたチャネルにおいて、" + "`" + config.CurrentConfig.Guild.Prefix + "` の後にコマンド名を打ってください。\n本 Bot インスタンスに指定されたチャネルは <#" + config.CurrentConfig.Guild.Channel + "> です。\n# コマンド一覧\n## `ping`\nsumika の生存確認をします。\n## `help`\nsumika の使い方を表示します。\n## `get`\n### Usage\n`" + config.CurrentConfig.Guild.Prefix + "get ${URL}`\nURL から動画を取得し、カラオケ用に最適化します。成功した場合、しばらくするとダウンロードできる URL を返信します。"
	msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
}

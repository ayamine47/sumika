package embed

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	ColorPrimary = 0xD0F3BA // #D0F3BA
	ColorSuccess = 0x587250 // #587250
	ColorError   = 0xE84C28 // #E84C28
)

type MsgInfo struct {
	Session *discordgo.Session
	OrgMsg  *discordgo.MessageCreate
}

func NewEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate) *discordgo.MessageEmbed {
	now := time.Now()
	msg := &discordgo.MessageEmbed{}
	msg.Author = &discordgo.MessageEmbedAuthor{}
	msg.Footer = &discordgo.MessageEmbedFooter{}
	msg.Author.IconURL = session.State.User.AvatarURL("256")
	msg.Author.Name = session.State.User.Username
	msg.Footer.IconURL = orgMsg.Author.AvatarURL("256")
	msg.Footer.Text = "Request from " + orgMsg.Author.Username + " @ " + now.String()
	msg.Color = ColorPrimary
	return msg
}

func NewSuccessEmbed(msgInfo *MsgInfo, description string) *discordgo.MessageEmbed {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Color = ColorSuccess
	msg.Title = "✅ Success"
	msg.Description = description
	return msg
}

func NewErrorEmbed(msgInfo *MsgInfo, description string) *discordgo.MessageEmbed {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Color = ColorError
	msg.Title = "❌ Error"
	msg.Description = description
	return msg
}

func SendMessageEmbed(msgInfo *MsgInfo, title string, description string) error {
	msg := NewEmbed(msgInfo.Session, msgInfo.OrgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(title)
	msg.Description = description
	_, err := msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, msg)
	return err
}

func SendSuccessEmbed(msgInfo *MsgInfo, description string, reference *discordgo.MessageReference) error {
	_, err := msgInfo.Session.ChannelMessageSendEmbedReply(msgInfo.OrgMsg.ChannelID, NewSuccessEmbed(msgInfo, description), reference)
	return err
}

func SendErrorEmbed(msgInfo *MsgInfo, description string) error {
	_, err := msgInfo.Session.ChannelMessageSendEmbed(msgInfo.OrgMsg.ChannelID, NewErrorEmbed(msgInfo, description))
	return err
}

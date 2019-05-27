package main

import (
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func Load() {

}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Avatar", "Provides a link to your avatar.", "!avatar", common.PRIORITY_LOW}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	url := message.Author.AvatarURL("")
	embed := common.NewEmbed().
		AddField("Avatar URL", url).
		SetThumbnail(url).MessageEmbed
	session.ChannelMessageSendEmbed(message.ChannelID, embed)
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"avatar")
}

func IsAdminOnly() bool {
	return true
}

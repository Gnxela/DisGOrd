package main

import (
	"fmt"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Avatar", "Provides a link to your avatar.", "!avatar", common.PRIORITY_LOW}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, message.Author.AvatarURL("")))
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"avatar")
}

func IsAdminOnly() bool {
	return true
}

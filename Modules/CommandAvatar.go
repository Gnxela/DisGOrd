package main

import (
	"fmt"
	"strings"

	"../Common"

	"github.com/bwmarrin/discordgo"
)

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, message.Author.AvatarURL("")))
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "avatar")
}

func IsAdminOnly() (bool) {
	return true
}

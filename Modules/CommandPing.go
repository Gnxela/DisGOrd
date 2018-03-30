package main

import (
	"strings"
	"fmt"

	"../Common"

	"github.com/bwmarrin/discordgo"
)

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, pong!", message.Author.ID))
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "ping")
}

func IsAdminOnly() (bool) {
	return true
}

package main

import (
	"fmt"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Ping", "Pongs and pings. Used for testing only.", "!ping", common.PRIORITY_HIGHEST}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, pong!", message.Author.ID))
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"ping")
}

func IsAdminOnly() bool {
	return true
}

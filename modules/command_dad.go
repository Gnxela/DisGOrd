package main

import (
	"fmt"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Dad", "Makes funny jokes.", "", common.PRIORITY_LOWEST}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	name := message.Content[strings.Index(message.Content, " ")+1:]
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hi %s, I'm Dad.", name))
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	m := strings.ToLower(message.Content)
	return strings.HasPrefix(m, "im ") || strings.HasPrefix(m, "i'm ")
}

func IsAdminOnly() bool {
	return false
}

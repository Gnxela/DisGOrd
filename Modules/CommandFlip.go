package main

import (
	"fmt"
	"math/rand"
	"strings"

	"../Common"

	"github.com/bwmarrin/discordgo"
)

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) (bool) {
	result := rand.Intn(2)
	if result == 0 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, heads!", message.Author.ID))
	} else {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, tails!", message.Author.ID))
	}
	return false
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.TrimRight(message.Content, "\n") ==  bot.Prefix+"flip"
}

func IsAdminOnly() bool {
	return false
}

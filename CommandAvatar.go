package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandAvatar struct {

}

func (c CommandAvatar) Fire(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, https://cdn.discordapp.com/avatars/%s/%s.png", message.Author.ID, message.Author.ID, message.Author.Avatar))
}

func (c CommandAvatar) ShouldFire(bot *Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "avatar")
}

func (c CommandAvatar) IsAdminOnly() (bool) {
	return true
}

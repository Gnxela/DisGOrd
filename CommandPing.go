package main

import (
	"strings"
		
	"github.com/bwmarrin/discordgo"
)

type CommandPing struct {
	
}

func (c CommandPing) Fire(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	session.ChannelMessageSend(message.ChannelID, "<@" + message.Author.ID + ">, pong!")
}

func (c CommandPing) ShouldFire(bot *Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "ping")
}

func (c CommandPing) IsAdminOnly() (bool) {
	return true
}
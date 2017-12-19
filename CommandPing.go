package main

import (
	"strings"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandPing struct {
	count int
}

func NewCommandPing() (c CommandPing) {
	ping := CommandPing{}
	ping.count = 0;
	return ping
}

func (c *CommandPing) Fire(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	c.count++
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, pong! (%d)", message.Author.ID, c.count))
}

func (c *CommandPing) ShouldFire(bot *Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "ping")
}

func (c *CommandPing) IsAdminOnly() (bool) {
	return true
}

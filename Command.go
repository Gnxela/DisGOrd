package main

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Fire(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate)
	ShouldFire(bot *Bot, message *discordgo.MessageCreate) bool
	IsAdminOnly() bool
}

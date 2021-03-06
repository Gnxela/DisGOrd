package common

import (
	"github.com/bwmarrin/discordgo"
)

type Module struct {
	Module      string
	Load        func()
	Unload      func()
	GetData     func(bot *Bot) Data
	Fire        func(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool
	ShouldFire  func(bot *Bot, message *discordgo.MessageCreate) bool
	IsAdminOnly func() bool
}

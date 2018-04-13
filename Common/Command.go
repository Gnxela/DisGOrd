package common

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	/*Fire(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate)
	ShouldFire(bot *Bot, message *discordgo.MessageCreate) bool
	IsAdminOnly() bool*/
	Module string
	Fire func(bot *Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool
	ShouldFire func(bot *Bot, message *discordgo.MessageCreate) bool
	IsAdminOnly func() bool
}

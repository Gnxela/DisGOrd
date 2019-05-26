package main

import (
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func Load() {

}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Ping", "Pongs and pings. Used for testing only.", "!ping", common.PRIORITY_HIGHEST}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	ping := session.HeartbeatLatency()
	embed := common.NewEmbed().
		AddField("Ping", ping.String()).
		SetColor(0x22aaff).InlineAllFields().MessageEmbed
	session.ChannelMessageSendEmbed(message.ChannelID, embed)
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"ping")
}

func IsAdminOnly() bool {
	return true
}

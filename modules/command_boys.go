package main

import (
	"fmt"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func Load() {

}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Boys", "@'s all online users (excluding away users) with the role \"Boys\".", "!boys", common.PRIORITY_MEDIUM}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	guildMapping := bot.ChannelMap[message.ChannelID]
	guild, err := session.Guild(guildMapping.ID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, err))
		return true
	}

	if guild.Unavailable {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, guild is unavailable. Please try again later.")
		return true
	}

	presenceSet := make(map[string]*discordgo.Presence, len(guild.Presences))
	for _, presence := range guild.Presences {
		presenceSet[presence.User.ID] = presence
	}
	result := "The following Boys have been summoned: "
	for _, member := range guild.Members {
		presence := presenceSet[member.User.ID]
		if presence != nil {
			if !member.User.Bot && presence.Status == discordgo.StatusOnline {
				for _, role := range member.Roles {
					if role == "435017423712813066" {
						result += fmt.Sprintf("<@%s> ", member.User.ID)
						break
					}
				}
			}
		}
	}
	session.ChannelMessageSend(message.ChannelID, result)
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"boys")
}

func IsAdminOnly() bool {
	return false
}

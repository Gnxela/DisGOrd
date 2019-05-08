package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func Load() {

}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Roll", "Rolls a dice with the number of sides specified.", "!roll <sides>", common.PRIORITY_MEDIUM}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	strs := strings.Split(message.Content, " ")
	if len(strs) < 2 {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !roll <sides>")
		return true
	}
	str := strs[1]
	sides, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !roll <sides>")
	} else {
		if sides <= 0 {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, '%d' in not a valid number.", message.Author.ID, sides))
		} else {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, rolling a dice with %s sides... %d", message.Author.ID, str, rand.Intn(int(sides))+1))
		}
	}
	return false
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"roll")
}

func IsAdminOnly() bool {
	return false
}

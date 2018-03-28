package main

import (
	"fmt"
	"strings"
	"strconv"
	"math/rand"

	"./Common"

	"github.com/bwmarrin/discordgo"
)

type CommandRoll struct {

}

func (c *CommandRoll) Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	strs := strings.Split(message.Content, " ")
	if len(strs) < 2 {
		session.ChannelMessageSend(message.ChannelID, "<@" + message.Author.ID + ">, !roll <sides>")
		return
	}
	str := strs[1]
	sides, err := strconv.ParseInt(str, 10, 32)
	if(err != nil) {
		session.ChannelMessageSend(message.ChannelID, "<@" + message.Author.ID + ">, !roll <sides>")
	} else {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, rolling a dice with %s sides... %d", message.Author.ID, str, rand.Intn(int(sides)) + 1))
	}
}

func (c *CommandRoll) ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) (bool) {
	return strings.HasPrefix(message.Content, bot.Prefix + "roll")
}

func (c *CommandRoll) IsAdminOnly() (bool) {
	return true
}

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Timer", "Sets a timer for the specified number of seconds. The caller will be @'d when the timer is finished.", "!timer <seconds>", common.PRIORITY_MEDIUM}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	strs := strings.Split(message.Content, " ")
	if len(strs) < 2 {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !timer <time (s)>")
		return true
	}
	str := strs[1]
	length, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !timer <time (s)>")
		return true
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, timer set.", message.Author.ID))
	go wait(length, session, message)
	return true
}

func wait(length int64, session *discordgo.Session, message *discordgo.MessageCreate) {
	<-time.After(time.Duration(length) * time.Second)
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, your timer is up!", message.Author.ID))
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"timer")
}

func IsAdminOnly() bool {
	return false
}

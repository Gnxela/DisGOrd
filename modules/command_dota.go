package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"../common"

	"github.com/bwmarrin/discordgo"
)

func Load() {

}

func Unload() {

}

var (
	active map[string]struct{} = make(map[string]struct{})
)

func GetData(bot *common.Bot) common.Data {
	return common.Data{"dota", "User will be @'d before every rune spawn until !stopdota is used (or after an hour).", "!dota <time as mm:ss, before first rune at 00:00> !stopdota", common.PRIORITY_MEDIUM}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	strs := strings.Split(message.Content, " ")
	if strs[0] == bot.Prefix+"dota" {
		if len(strs) < 2 {
			session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !dota <time as mm:ss, before first rune at 00:00>")
			return true
		}
		timestamp := strs[1]
		var start int64
		if strings.Contains(timestamp, ":") {
			splitstamp := strings.Split(timestamp, ":")
			minutes, err := strconv.ParseInt(splitstamp[0], 10, 32)
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !dota <time as mm:ss, before first rune at 00:00>")
				return true
			}
			seconds, err := strconv.ParseInt(splitstamp[1], 10, 32)
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !dota <time as mm:ss, before first rune at 00:00>")
				return true
			}
			start = minutes*60 + seconds
		} else {
			val, err := strconv.ParseInt(timestamp, 10, 32)
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !dota <time as mm:ss, before first rune at 00:00>")
				return true
			}
			start = val
		}
		go run(start, session, message)
	} else {
		delete(active, message.Author.ID)
	}
	return true
}

func run(start int64, session *discordgo.Session, message *discordgo.MessageCreate) {
	i := 0
	var diff int64
	active[message.Author.ID] = struct{}{}
	if start > 15 {
		<-time.After(time.Duration(start-15) * time.Second)
	} else {
		diff = 15 - start
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, runes!", message.Author.ID))
	for i <= 18 {
		<-time.After(5*time.Minute - time.Duration(diff)*time.Second)
		diff = 0
		_, ok := active[message.Author.ID]
		if !ok {
			break
		}
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, runes!", message.Author.ID))
		i++
	}
	delete(active, message.Author.ID)
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"dota") || strings.HasPrefix(message.Content, bot.Prefix+"stopdota")
}

func IsAdminOnly() bool {
	return false
}

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

var (
	lexer *common.Lexer = common.CreateLexer(common.CreateSequence(&common.AbsoluteToken{"!in"}, &common.NumericalToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!in"}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!in"}),
	)
)

func Load() {

}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"In", "@'s all users in a voice channel.", "!in [channel name]. !in [channelID]", common.PRIORITY_LOW}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	i, values := lexer.ParseCommand(message.Content)
	switch i {
	case 0: //ID
		err := callChannel(message.ChannelID, strconv.Itoa(int(values[1].(int64))), message.GuildID, session)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s.", message.Author.ID, err))
		}
	case 1: //Name
		channelID, err := findChannel(session, message.GuildID, values[1].(string))
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s.", message.Author.ID, err))
		}
		err = callChannel(message.ChannelID, channelID, message.GuildID, session)
		if err != nil {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s.", message.Author.ID, err))
		}
	default: //Call role
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s.", message.Author.ID, GetData(bot).Usage))
	}
	return true
}

func findChannel(session *discordgo.Session, guildID, name string) (string, error) {
	guild, err := session.Guild(guildID)
	if err != nil {
		return "", err
	}
	name = strings.ToLower(name)
	for _, channel := range guild.Channels {
		if strings.ToLower(channel.Name) == name {
			return channel.ID, nil
		}
	}
	return "", errors.New("channel not found")
}

func callChannel(textChannelID, voiceChannelID, guildID string, session *discordgo.Session) error {
	guild, err := session.Guild(guildID)
	if err != nil {
		return err
	}
	str := "The following users have been called: "
	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == voiceChannelID {
			str += fmt.Sprintf("<@%s> ", voiceState.UserID)
		}
	}
	session.ChannelMessageSend(textChannelID, str)
	return nil
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, "!in")
}

func IsAdminOnly() bool {
	return false
}

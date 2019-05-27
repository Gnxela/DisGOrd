package main

import (
	"fmt"

	"../common"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Roles map[string]string //Map from alias to RoleID
}

var (
	configFile string = "config/command_call.json"
	config     Config
	lexer      *common.Lexer
)

func createLexer() {
	sequences := []*common.Sequence{common.CreateSequence(&common.AbsoluteToken{"!call"}, &common.AbsoluteToken{"add"}, &common.StringToken{}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!call"}, &common.AbsoluteToken{"remove"}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!call"}, &common.AbsoluteToken{"list"}),
		common.CreateSequence(&common.AbsoluteToken{"!call"}),
	}
	for alias, _ := range config.Roles {
		sequences = append(sequences, common.CreateSequence(&common.AbsoluteToken{"!" + alias}))
	}
	lexer = common.CreateLexer(sequences...)
}

func Load() {
	err := common.LoadConfig(configFile, &config)
	if err != nil {
		config = Config{make(map[string]string, 0)}
		common.SaveConfig(configFile, config)
	}
	createLexer()
}

func Unload() {
	common.SaveConfig(configFile, config)
}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Call", "@'s all online users with specific roles.", "!call list. !call add [alias] [RollID]. !call remove [alias]. ![alias]", common.PRIORITY_LOW}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	i, values := lexer.ParseCommand(message.Content)
	switch i {
	case 0: //Add
		config.Roles[values[2].(string)] = values[3].(string)
		createLexer()
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, role added.", message.Author.ID))
	case 1: //Remove
		delete(config.Roles, values[2].(string))
		createLexer()
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, role removed.", message.Author.ID))
	case 2: //List
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %v", message.Author.ID, config.Roles))
	case 3: //Invalid
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, GetData(bot).Usage))
	default: //Call role
		alias := message.Content[1:]
		roleID := config.Roles[alias]
		callRole(roleID, bot, session, message)
	}
	return true
}

func callRole(roleID string, bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	guild, err := session.Guild(message.GuildID)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, err))
		return
	}

	if guild.Unavailable {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, guild is unavailable. Please try again later.")
		return
	}

	//Map from UserID to Presence. Map is reduced until only targeted users are contained.
	presenceSet := make(map[string]*discordgo.Presence, len(guild.Presences))
	for _, presence := range guild.Presences {
		presenceSet[presence.User.ID] = presence
	}
	result := "The following users have been called: "
L:
	for _, member := range guild.Members {
		presence := presenceSet[member.User.ID]
		if presence == nil || member.User.Bot || presence.Status != discordgo.StatusOnline {
			delete(presenceSet, member.User.ID)
			continue
		}
		for _, role := range member.Roles {
			if role == roleID {
				continue L //Found the role, don't remove this user.
			}
		}
		delete(presenceSet, member.User.ID)
	}
	for id, _ := range presenceSet {
		result += "<@" + id + "> "
	}
	session.ChannelMessageSend(message.ChannelID, result)
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	i, _ := lexer.ParseCommand(message.Content)
	return i >= 0
}

func IsAdminOnly() bool {
	return false
}

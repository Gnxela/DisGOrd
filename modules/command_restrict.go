package main

import (
	"fmt"
	"strings"

	"../common"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Enabled   bool
	Whitelist bool
	Channels  map[string]struct{}
}

var (
	//This does not allow for the changing of the prefix.
	lexer *common.Lexer = common.CreateLexer(common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"enable"}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"disable"}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"add"}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"remove"}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"blacklist"}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"whitelist"}),
	)
	config     Config
	configFile string = "config/command_restrict.json"
)

func Load() {
	err := common.LoadConfig(configFile, &config)
	if err != nil {
		config = Config{false, false, make(map[string]struct{}, 0)}
		common.SaveConfig(configFile, config)
	}
}

func Unload() {
	common.SaveConfig(configFile, config)
}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Restrict", "Restricts the channels in which commands can be used.", "!restrict {enable|disable}. !restrict {add|remove} [ChannelID]. !restrict {blacklist|whitelist}", common.PRIORITY_HIGHEST}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if strings.HasPrefix(message.Content, bot.Prefix+"restrict") {
		i, values := lexer.ParseCommand(message.Content)
		switch i {
		case 0: //Enable
			config.Enabled = true
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, enabled.", message.Author.ID))
		case 1: //Disable
			config.Enabled = false
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, disabled.", message.Author.ID))
		case 2: //Add
			config.Channels[values[2].(string)] = struct{}{}
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, channel added.", message.Author.ID))
		case 3: //Remove
			delete(config.Channels, values[2].(string))
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, channel removed.", message.Author.ID))
		case 4: //Blacklist
			config.Whitelist = false
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, blacklist enabled.", message.Author.ID))
		case 5: //Whitelist
			config.Whitelist = true
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, whitelist enabled.", message.Author.ID))
		default:
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s.", message.Author.ID, GetData(bot).Usage))
			return true
		}
		common.SaveConfig(configFile, config)
		return true
	}
	if config.Enabled {
		_, ok := config.Channels[message.ChannelID]
		return ok == config.Whitelist
	} else {
		return true
	}
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return true
}

func IsAdminOnly() bool {
	return true
}

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
	lexer *common.Lexer = common.CreateLexer(common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"whitelist"}, &common.StringToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!restrict"}, &common.AbsoluteToken{"blacklist"}, &common.StringToken{}),
	)
	dummy      int = Load()
	config     *Config
	configFile string = "config/command_restrict.json"
)

func Load() int {
	err := common.LoadConfig(configFile, config)
	if err != nil {
		config = &Config{false, false, make(map[string]struct{}, 0)}
		common.SaveConfig(configFile, config)
	}
	return 0
}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Restrict", "Restricts the channels in which commands can be used.", "!restrict {enable|disable|whitelist|blacklist} [channelID]", common.PRIORITY_HIGHEST}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if strings.HasPrefix(message.Content, bot.Prefix+"restrict") {
		i, values := lexer.ParseCommand(message.Content)
		if values == nil {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, GetData(bot).Usage))
		} else if i == 0 {
			fmt.Printf("%v\n", values)
		} else if i == 1 {
			fmt.Printf("%v\n", values)
		}
	}
	return false
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	return strings.HasPrefix(message.Content, bot.Prefix+"restrict")
}

func IsAdminOnly() bool {
	return true
}

package main

import (
	"fmt"
	"math/rand"

	"../common"

	"github.com/bwmarrin/discordgo"
)

var (
	lexer *common.Lexer
)

func Load() {
	lexer = common.CreateLexer(common.CreateSequence(&common.AbsoluteToken{"!roll"}, &common.NumericalToken{}),
		common.CreateSequence(&common.AbsoluteToken{"!roll"}),
	)
}

func Unload() {

}

func GetData(bot *common.Bot) common.Data {
	return common.Data{"Roll", "Rolls a dice with the number of sides specified.", "!roll <sides>", common.PRIORITY_MEDIUM}
}

func Fire(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) bool {
	i, values := lexer.ParseCommand(message.Content)
	if i == 1 {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !roll <sides>")
		return true
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, rolling a dice with %d sides... %d", message.Author.ID, values[1].(int64), rand.Intn(int(values[1].(int64)))+1))
	return true
}

func ShouldFire(bot *common.Bot, message *discordgo.MessageCreate) bool {
	i, _ := lexer.ParseCommand(message.Content)
	return i >= 0
}

func IsAdminOnly() bool {
	return false
}

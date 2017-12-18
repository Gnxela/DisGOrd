package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Prefix string
	Commands []Command
}

var bot Bot = Bot{"!", make([]Command, 0)};
var token string;

func init() {
	flag.StringVar(&token, "t", "", "Token")
	flag.Parse()
	if token == "" {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	bot.Commands = append(bot.Commands, CommandPing{}, CommandRoll{}, CommandAvatar{})

	discord.AddHandler(onMessage)

	err = discord.Open();
	if err != nil {
		fmt.Println("Error creating WebSocket,", err)
		return
	}

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-close

	discord.Close()
}

func onMessage(session *discordgo.Session, mesage *discordgo.MessageCreate) {
	if mesage.Author.ID == session.State.User.ID {
		return
	}

	for _, element := range bot.Commands {
		if(element.ShouldFire(&bot, mesage)) {
			element.Fire(&bot, session, mesage)
		}
	}
}
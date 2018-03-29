package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"plugin"
	"strings"
	"io/ioutil"
	"bytes"

	"./Common"

	"github.com/bwmarrin/discordgo"
)

var bot common.Bot = common.Bot{"!", make([]common.Command, 0)}
var token string

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

	initCommands()

	discord.AddHandler(onMessage)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error creating WebSocket,", err)
		return
	}

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-close

	discord.Close()
}

func initCommands() {

}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {//Ignore messages from ourselves
		return
	}

	//Three designated commands.
	if strings.HasPrefix(message.Content, bot.Prefix + "load") {
		load(&bot, session, message)
	} else if strings.HasPrefix(message.Content, bot.Prefix + "unload") {
		
	} else if strings.HasPrefix(message.Content, bot.Prefix + "list") {
		list(&bot, session, message)
	} else {
		for _, element := range bot.Commands {
			if(element.ShouldFire(&bot, message)) {
				element.Fire(&bot, session, message)
			}
		}
	}
}

func load(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	strs := strings.SplitN(message.Content, " ", 2)
	if len(strs) != 2 {
		session.ChannelMessageSend(message.ChannelID, "<@" + message.Author.ID + ">, !load <module>")
		return
	}
	module := strs[1]
	files, err := ioutil.ReadDir("./Modules")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.EqualFold(module, file.Name()) {
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, loading '%s'.", message.Author.ID, module))
			p, err := plugin.Open("./Modules/" + module)
			if err != nil {
				panic(err)
			}
			command := common.Command{}
			fire, err := p.Lookup("Fire")
			if err != nil {
				panic(err)
			}
			shouldFire, err := p.Lookup("ShouldFire")
			if err != nil {
				panic(err)
			}
			isAdminOnly, err := p.Lookup("IsAdminOnly")
			if err != nil {
				panic(err)
			}
			command.Module = module
			command.Fire = fire.(func(*common.Bot, *discordgo.Session, *discordgo.MessageCreate))
			command.ShouldFire = shouldFire.(func(*common.Bot, *discordgo.MessageCreate) bool)
			command.IsAdminOnly = isAdminOnly.(func() bool)
			bot.Commands = append(bot.Commands, command)
		}
	}
}

func list(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	files, err := ioutil.ReadDir("./Modules")
	if err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	buffer.WriteString("Modules:\n```diff\n")
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".so") {
			continue
		}
		found := false
		for _, module := range bot.Commands {
			if module.Module == file.Name() {
				found = true
				break
			}
		}
		if found {
			buffer.WriteString(fmt.Sprintf("+ %s\n", file.Name()))
		} else {
			buffer.WriteString(fmt.Sprintf("- %s\n", file.Name()))
		}
	}
	buffer.WriteString("```")
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, %s", message.Author.ID, buffer.String()))
}

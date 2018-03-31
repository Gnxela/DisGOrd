package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"plugin"
	"strings"
	"io"
	"io/ioutil"
	"bytes"
	"encoding/json"

	"./Common"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	LoadedModules map[string]bool//The value will never be used.
}

var bot common.Bot = common.Bot{"!", make([]common.Command, 0)}
var config Config//Would store in bot, but don't think modules need access to it.
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
	loadConfig()

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	discord.AddHandler(onMessage)

	enableLoadedModules()

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

func enableLoadedModules() {
	for module, _ := range config.LoadedModules {
		loadModule(module)
	}
}

func saveConfig() {
	configFile, _ := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC, 0755)//Need to look into FileModes and general UNIX file permissions.
	defer configFile.Close()
	decoder := json.NewEncoder(configFile)
	err := decoder.Encode(&config)
	if err != nil {
		panic(err)
	}
}

func loadConfig() {
	configFile, _ := os.OpenFile("config.json", os.O_RDONLY|os.O_CREATE, 0755)//Need to look into FileModes and general UNIX file permissions.
	defer configFile.Close()
	decoder := json.NewDecoder(configFile)
	config = Config{}
	err := decoder.Decode(&config)
	if err != nil {
		if err == io.EOF {
			config.LoadedModules = make(map[string]bool, 0)//When no file is read, map is never initialised, so we need to do it manually.
			saveConfig()
			return
		}
		panic(err)
	}
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {//Ignore messages from ourselves
		return
	}

	//Three designated commands.
	if strings.HasPrefix(message.Content, bot.Prefix + "load") {
		load(&bot, session, message)
	} else if strings.HasPrefix(message.Content, bot.Prefix + "unload") {
		unload(&bot, session, message)
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
			err := loadModule(module)
			if err != nil {
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, failed to load '%s'. (%s)", message.Author.ID, module, err))
			} else {
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, loaded '%s'.", message.Author.ID, module))
			}
		}
	}
}

func unload(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	strs := strings.SplitN(message.Content, " ", 2)
	if len(strs) != 2 {
		session.ChannelMessageSend(message.ChannelID, "<@" + message.Author.ID + ">, !unload <module>")
		return
	}
	module := strs[1]
	for index, mod := range bot.Commands {
		if mod.Module == module {
			/*
				According to the docs: "A plugin is only initialized once, and cannot be closed.A plugin is only initialized once, and cannot be closed."
				I need to look into what the garbaage disposal of Go, and see if me reloading the plugins is bad.
				If plugins can't be unloaded and cleaned by GC, then I need to store them and reuse them, to avoid memory problems.
			*/
			bot.Commands = append(bot.Commands[:index], bot.Commands[index + 1:]...)
			delete(config.LoadedModules, module)
			saveConfig()
			session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, unloaded '%s'.", message.Author.ID, module))
			return
		}
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, failed to unload '%s'.", message.Author.ID, module))
}

func loadModule(module string) (err error) {
	p, err := plugin.Open("./Modules/" + module)
	if err != nil {
		return
	}
	command := common.Command{}
	fire, err := p.Lookup("Fire")
	if err != nil {
		return
	}
	shouldFire, err := p.Lookup("ShouldFire")
	if err != nil {
		return
	}
	isAdminOnly, err := p.Lookup("IsAdminOnly")
	if err != nil {
		return
	}
	command.Module = module
	command.Fire = fire.(func(*common.Bot, *discordgo.Session, *discordgo.MessageCreate))
	command.ShouldFire = shouldFire.(func(*common.Bot, *discordgo.MessageCreate) bool)
	command.IsAdminOnly = isAdminOnly.(func() bool)
	bot.Commands = append(bot.Commands, command)
	config.LoadedModules[module] = true
	saveConfig()
	fmt.Printf("Loaded %s\n", module)
	return
}

func list(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	files, err := ioutil.ReadDir("./Modules")
	if err != nil {
		panic(err)
	}

	var loadedBuffer bytes.Buffer
	var unloadedBuffer bytes.Buffer
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
			loadedBuffer.WriteString(fmt.Sprintf("+ %s\n", file.Name()))
		} else {
			unloadedBuffer.WriteString(fmt.Sprintf("- %s\n", file.Name()))
		}
	}
	if loadedBuffer.String() == "" {
		loadedBuffer.WriteString("  None\n")
	}
	if unloadedBuffer.String() == "" {
		unloadedBuffer.WriteString("  None\n")
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, ```diff\nLoaded Modules:\n%sUnloaded Modules:\n%s```", message.Author.ID, loadedBuffer.String(), unloadedBuffer.String()))
}

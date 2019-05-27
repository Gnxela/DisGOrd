package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"plugin"
	"strings"
	"syscall"

	"./common"

	"github.com/bwmarrin/discordgo"
)

// Config for the bot that is loaded from a .json
type Config struct {
	Token         string
	LoadedModules map[string]struct{}
	Bin           string
}

var (
	bot        *common.Bot = &common.Bot{make(map[string]*common.Guild, 0), "!", make(map[common.Priority][]*common.Module, 0)}
	config     Config
	configFile string = "config/config.json"
)

func init() {
	common.LoadConfig(configFile, &config)
	if config.LoadedModules == nil {
		config.LoadedModules = make(map[string]struct{}, 0)
	}
	if config.Token == "" { //Load token from command line
		t := ""
		flag.StringVar(&t, "t", "", "Token")
		flag.Parse()
		if t == "" {
			flag.Usage()
			os.Exit(0)
		}
		config.Token = t
		common.SaveConfig(configFile, config)
	}
}

func main() {
	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Error creating Discord session,", err)
		return
	}

	discord.AddHandler(onMessage)
	discord.AddHandlerOnce(onReady)

	enableLoadedModules()

	err = discord.Open()
	defer discord.Close()
	if err != nil {
		fmt.Println("Error creating WebSocket,", err)
		return
	}

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-close

	unloadModules()
}

func unloadModules() {
	for _, modules := range bot.Modules {
		for _, module := range modules {
			module.Unload()
		}
	}
}

func enableLoadedModules() {
	failed := false
	for module := range config.LoadedModules {
		err := loadModule(module)
		if err != nil {
			fmt.Printf("Failed to load '%s': %s\n", module, err)
			delete(config.LoadedModules, module)
			failed = true
		}
	}
	if failed {
		common.SaveConfig(configFile, config)
	}
}

func onReady(session *discordgo.Session, ready *discordgo.Ready) {
	guilds, err := session.UserGuilds(100, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Logged into %d guilds.\n", len(guilds))
	for _, g := range guilds {
		guild := &common.Guild{g.ID, false}
		bot.Guilds[g.ID] = guild
		go loadGuild(session, guild)
	}
}

func loadGuild(session *discordgo.Session, guild *common.Guild) {
	_, err := session.Guild(guild.ID)
	if err != nil {
		panic(err)
	}

	guild.Ready = true
	fmt.Printf("Loaded guild: %s\n", guild.ID)
}

func onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !bot.Guilds[message.GuildID].Ready { //Ignore messages to guilds that aren't ready.
		return
	}

	if message.Author.Bot { //Ignore messages from bots.
		return
	}

	if message.Author.ID == session.State.User.ID { //Ignore messages from ourselves.
		return
	}

	//Three designated commands. Lazy evaluation means we wont be checking for admin unnessessarily.
	if strings.HasPrefix(message.Content, bot.Prefix+"load") && common.CheckAdmin(session, message.ChannelID, message.Author.ID) {
		load(bot, session, message)
	} else if strings.HasPrefix(message.Content, bot.Prefix+"unload") && common.CheckAdmin(session, message.ChannelID, message.Author.ID) {
		unload(bot, session, message)
	} else if strings.HasPrefix(message.Content, bot.Prefix+"list") && common.CheckAdmin(session, message.ChannelID, message.Author.ID) {
		list(bot, session, message)
	} else {
	L:
		for priority := int(common.PRIORITY_CANCEL); priority <= int(common.PRIORITY_OBSERVE); priority++ {
			modules := bot.Modules[common.Priority(priority)]
			for _, module := range modules {
				if !module.IsAdminOnly() || (module.IsAdminOnly() && common.CheckAdmin(session, message.ChannelID, message.Author.ID)) {
					if module.ShouldFire(bot, message) {
						if !module.Fire(bot, session, message) {
							break L
						}
					}
				}
			}
		}
	}
}

func load(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	strs := strings.SplitN(message.Content, " ", 2)
	if len(strs) != 2 {
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !load <module>")
		return
	}
	module := strs[1]
	files, err := ioutil.ReadDir(config.Bin)
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
		session.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+">, !unload <module>")
		return
	}
	module := strs[1]
	for priority, modules := range bot.Modules {
		for index, mod := range modules {
			if mod.Module == module {
				/*
					According to the docs: "A plugin is only initialized once, and cannot be closed.A plugin is only initialized once, and cannot be closed."
					I need to look into what the garbaage disposal of Go, and see if me reloading the plugins is bad.
					If plugins can't be unloaded and cleaned by GC, then I need to store them and reuse them, to avoid memory problems.
				*/
				mod.Unload()
				bot.Modules[priority] = append(bot.Modules[priority][:index], bot.Modules[priority][index+1:]...)
				delete(config.LoadedModules, module)
				common.SaveConfig(configFile, config)
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, unloaded '%s'.", message.Author.ID, module))
				return
			}
		}
	}
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s>, failed to unload '%s'. Module not found (maybe it wasn't loaded?).", message.Author.ID, module))
}

func loadModule(moduleName string) (err error) {
	p, err := plugin.Open(config.Bin + moduleName)
	if err != nil {
		return
	}
	module := common.Module{}
	load, err := p.Lookup("Load")
	if err != nil {
		return
	}
	unload, err := p.Lookup("Unload")
	if err != nil {
		return
	}
	getData, err := p.Lookup("GetData")
	if err != nil {
		return
	}
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%s failed to load: %s\n", module, r)
			switch t := r.(type) {
			case error:
				err = t
			case string:
				err = errors.New(t)
			default:
				err = errors.New("loadModule() unknown recovery")
			}
		}
	}()
	module.Module = moduleName
	module.Load = load.(func())
	module.Unload = unload.(func())
	module.GetData = getData.(func(*common.Bot) common.Data)
	module.Fire = fire.(func(*common.Bot, *discordgo.Session, *discordgo.MessageCreate) bool)
	module.ShouldFire = shouldFire.(func(*common.Bot, *discordgo.MessageCreate) bool)
	module.IsAdminOnly = isAdminOnly.(func() bool)

	data := module.GetData(bot)

	bot.Modules[data.Priority] = append(bot.Modules[data.Priority], &module)
	module.Load()
	config.LoadedModules[moduleName] = struct{}{}
	common.SaveConfig(configFile, config)
	fmt.Printf("Loaded %s.\n", moduleName)
	return
}

func list(bot *common.Bot, session *discordgo.Session, message *discordgo.MessageCreate) {
	files, err := ioutil.ReadDir(config.Bin)
	if err != nil {
		panic(err)
	}

	var loadedBuffer bytes.Buffer
	var unloadedBuffer bytes.Buffer
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".so") {
			continue
		}
		found := 0
		for _, modules := range bot.Modules {
			for _, module := range modules {
				if module.Module == file.Name() {
					found++
				}
			}
		}
		if found == 1 {
			loadedBuffer.WriteString(file.Name())
			loadedBuffer.WriteString("\n")
		} else if found > 1 {
			loadedBuffer.WriteString(fmt.Sprintf("%s (%d) [Multible instances of plugin loaded]\n", file.Name(), found))
		} else {
			unloadedBuffer.WriteString(file.Name())
			unloadedBuffer.WriteString("\n")
		}
	}
	if loadedBuffer.String() == "" {
		loadedBuffer.WriteString("  None\n")
	}
	if unloadedBuffer.String() == "" {
		unloadedBuffer.WriteString("  None\n")
	}
	embed := common.NewEmbed().
		AddField("Loaded", loadedBuffer.String()).
		AddField("Unloaded", unloadedBuffer.String()).
		InlineAllFields().MessageEmbed
	session.ChannelMessageSendEmbed(message.ChannelID, embed)
}

package common

import "github.com/bwmarrin/discordgo"

type Bot struct {
	Guilds []*Guild
	ChannelMap map[string]*Guild//A map that maps channel IDs to their guild. Populated when a guild is loaded
	Prefix string
	Commands []*Command
}

type Guild struct {
	Guild *discordgo.UserGuild
	Ready bool
}

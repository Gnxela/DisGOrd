package common

type Bot struct {
	Guilds     []*Guild
	ChannelMap map[string]*Guild //A map that maps channel IDs to their guild. Populated when a guild is loaded
	Prefix     string
	Commands   map[Priority][]*Command
}

type Guild struct {
	ID    string
	Ready bool
}

package common

type Bot struct {
	Guilds  map[string]*Guild
	Prefix  string
	Modules map[Priority][]*Module
}

type Guild struct {
	ID    string
	Ready bool
}

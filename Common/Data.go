package common

type Data struct {
	Name        string
	Description string
	Usage       string
	Priority    Priority //Only read once when a module is loaded.
}

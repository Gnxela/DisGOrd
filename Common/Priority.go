package common

type Priority int8

const (
	PRIORITY_CANCEL Priority = iota //Should only be used to cancel commands.
	PRIORITY_HIGHEST
	PRIORITY_HIGH
	PRIORITY_MEDIUM
	PRIORITY_LOW
	PRIORITY_LOWEST
	PRIORITY_OBSERVE //Should only be usued to observe commands.
)

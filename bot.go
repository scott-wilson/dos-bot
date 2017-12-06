package dosbot

import (
	"regexp"
)

type Bot interface {
	Name() string
	ID() interface{}
	DirectedMessageRegex() *regexp.Regexp
	SendMessage(Room, string) error
	SendDirectMessage(Room, User, string) error
	SendEmote(Room, string) error
	SendPrivateMessage(User, string) error
	SendPrivateEmote(User, string) error
}

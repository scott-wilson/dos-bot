package dosbot

import (
	"regexp"
)

type Bot interface {
	Name() string
	ID() interface{}
	DirectedMessageRegex() *regexp.Regexp
}

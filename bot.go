package dosbot

type Bot struct {
	name string
}

func NewBot(name string) (Bot, error) {
	return Bot{name: name}, nil
}

func (bot Bot) Name() string {
	return bot.name
}

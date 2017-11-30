package dosbot

const (
	EventDirectedMessage = "directed-message"
	EventChannelMessage  = "channel-message"
	EventTick            = "tick"
)

type Event struct {
	eventType string
	bot       Bot
	msg       string
	err       error
	sender    User
	room      Room
}

func NewEvent(eventType string, msg string, err error, sender User, room Room, bot Bot) Event {
	return Event{eventType: eventType, bot: bot, msg: msg, err: err, sender: sender, room: room}
}

func (e Event) Type() string {
	return e.eventType
}

func (e Event) Bot() Bot {
	return e.bot
}

func (e Event) Message() string {
	return e.msg
}

func (e Event) Error() error {
	return e.err
}

func (e Event) Sender() User {
	return e.sender
}

func (e Event) Room() Room {
	return e.room
}

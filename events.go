package dosbot

type Event struct {
	eventType string
	msg       string
	err       error
	sender    User
	room      Room
}

func NewEvent(eventType string, msg string, err error, sender User, room Room) Event {
	return Event{eventType: eventType, msg: msg, err: err, sender: sender, room: room}
}

func (e Event) Type() string {
	return e.eventType
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

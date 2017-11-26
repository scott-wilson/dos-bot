package dosbot

type Event struct {
	eventType string
	msg       string
	err       error
}

func NewEvent(eventType string, msg string, err error) Event {
	return Event{eventType: eventType, msg: msg, err: err}
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

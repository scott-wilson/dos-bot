package dosbot

import (
	"fmt"
	"log"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var registeredActions = make(map[string][]action)

var helpRegex = regexp.MustCompile(`(?i)help`)

func RegisterAction(eventName string, event func(Event) error, signature string, description string) {
	log.Printf("Registering action.\n\tEvent Name: %s\n\tEvent: %v\n\tSignature: %s\n\tDescription: %s\n", eventName, event, signature, description)
	action := action{Event: event, Signature: signature, Description: description}
	registeredActions[eventName] = append(registeredActions[eventName], action)
}

func EmitActions(event string, message string, sender User, room Room, bot Bot, connector chan<- Event) {
	log.Printf("Emitting action.\n\tEvent: %s\n\tMessage: %s\n\tUser: %#v\n\tRoom: %#v\n\tBot: %#v\n", event, message, sender, room, bot)
	connector <- NewEvent(event, message, nil, sender, room, bot)
}

func EmitDirectedMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	log.Printf("Emitting directed action.\n\tMessage: %s\n\tUser: %#v\n\tRoom: %#v\n\tBot: %#v\n", message, sender, room, bot)
	EmitActions(EventDirectedMessage, message, sender, room, bot, connector)
}

func EmitChannelMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	log.Printf("Emitting channel message action.\n\tMessage: %s\n\tUser: %#v\n\tRoom: %#v\n\tBot: %#v\n", message, sender, room, bot)
	EmitActions(EventChannelMessage, message, sender, room, bot, connector)
}

func EmitMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	log.Printf("Emitting message action.\n\tMessage: %s\n\tUser: %#v\n\tRoom: %#v\n\tBot: %#v\n", message, sender, room, bot)
	regex := bot.DirectedMessageRegex()
	result := regex.FindStringSubmatch(message)

	if len(result) > 0 {
		EmitDirectedMessageActions(result[1], sender, room, bot, connector)
	} else {
		EmitChannelMessageActions(message, sender, room, bot, connector)
	}
}

func EmitTickActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	EmitActions(EventTick, message, sender, room, bot, connector)
}

func actionWorker(toActions <-chan Event, toConnector chan<- Event) {
	for event := range toActions {
		actions, ok := registeredActions[event.Type()]

		if !ok {
			continue
		}

		for _, actionEvent := range actions {
			err := actionEvent.Event(event)
			toConnector <- NewEvent(event.Type(), "", err, event.Sender(), event.Room(), event.Bot())
		}
	}
}

func actionThreadPool(toActions <-chan Event, toConnector chan<- Event) {
	for w := 0; w < runtime.NumCPU(); w++ {
		go actionWorker(toActions, toConnector)
	}
}

type action struct {
	Event       func(Event) error
	Signature   string
	Description string
}

// ----------------
// Built in actions
// ----------------

func helpAction(event Event) error {
	if helpRegex.FindString(event.Message()) == "" {
		return nil
	}

	message := helpOnAllActions()
	sender := event.Sender()
	room := event.Room()
	bot := event.Bot()
	return bot.SendDirectMessage(room, sender, message)
}

func listAllActions() []action {
	allActions := make([]action, 0)
	for _, actions := range registeredActions {
		allActions = append(allActions, actions...)
	}

	return allActions
}

func helpOnAllActions() string {
	allActions := listAllActions()
	message := "Here are the actions that I understand:\n"

	sort.Slice(allActions, func(i, j int) bool {
		return allActions[i].Signature < allActions[j].Signature
	})

	allActionsHelp := make([]string, len(allActions))

	for index, action := range allActions {
		allActionsHelp[index] = fmt.Sprintf("%s: %s", action.Signature, action.Description)
	}

	actionList := strings.Join(allActionsHelp, "\n")
	message += actionList

	return message
}

func init() {
	RegisterAction(EventDirectedMessage, helpAction, "help", "Show list of supported actions with ")
}

package dosbot

import (
	"fmt"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var registeredActions = make(map[string][]action)

var helpRegex = regexp.MustCompile(`(?i)help[\t ]*(.*)`)

func RegisterAction(eventName string, event func(Event) error, title string, shortDescription string, longDescription string) {
	action := action{Event: event, Title: title, ShortDescription: shortDescription, LongDescription: longDescription}
	registeredActions[eventName] = append(registeredActions[eventName], action)
}

func EmitActions(event string, message string, sender User, room Room, bot Bot, connector chan<- Event) {
	connector <- NewEvent(event, message, nil, sender, room, bot)
}

func EmitDirectedMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	EmitActions(EventDirectedMessage, message, sender, room, bot, connector)
}

func EmitChannelMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
	EmitActions(EventChannelMessage, message, sender, room, bot, connector)
}

func EmitMessageActions(message string, sender User, room Room, bot Bot, connector chan<- Event) {
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
	Event            func(Event) error
	Title            string
	ShortDescription string
	LongDescription  string
}

// ----------------
// Built in actions
// ----------------

func helpAction(event Event) error {
	message := ""

	result := helpRegex.FindStringSubmatch(event.Message())

	if result[1] != "" {
		message = "Could not find that action."
		toSearch := strings.ToLower(result[1])

		for _, action := range listAllActions() {
			if strings.ToLower(action.Title) == toSearch {
				message = helpOnAction(action)
				break
			}
		}
	} else {
		message = helpOnAllActions()
	}

	if message == "" {
		return nil
	}

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
		return allActions[i].Title < allActions[j].Title
	})

	allActionsHelp := make([]string, len(allActions))

	for index, action := range allActions {
		allActionsHelp[index] = fmt.Sprintf("%s - %s", action.Title, action.ShortDescription)
	}

	actionList := strings.Join(allActionsHelp, "\n")
	message += actionList

	return message
}

func helpOnAction(action action) string {
	return fmt.Sprintf("%s\n\n%s\n\n%s", action.Title, action.ShortDescription, action.LongDescription)
}

func init() {
	RegisterAction(EventDirectedMessage, helpAction, "Help", "Show list of supported actions with @dos help or more information about a specific action with @dos help <action>", "Nothing more to say about that...")
}

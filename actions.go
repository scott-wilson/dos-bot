package dosbot

import (
	"errors"
	"runtime"
)

var (
	ErrInvalidEvent              = errors.New("event is invalid")
	ErrEventNotSupportedByAction = errors.New("action does not support event")
	ErrNoActionFound             = errors.New("no action found for event")
)

var registeredActions = make(map[string][]func(string) (string, error))

func RegisterAction(eventName string, action func(string) (string, error)) {
	registeredActions[eventName] = append(registeredActions[eventName], action)
}

func EmitActions(event string, message string, connector chan<- Event) {
	connector <- NewEvent(event, message, nil)
}

func actionWorker(eventsToActions <-chan Event, eventsToConnector chan<- Event) {
	for event := range eventsToActions {
		actions, ok := registeredActions[event.Type()]
		actionFound := false

		if !ok {
			continue
		}

		for _, action := range actions {
			result, err := action(event.Message())

			if err == ErrEventNotSupportedByAction {
				continue
			}

			actionFound = true
			eventsToConnector <- NewEvent(event.Type(), result, err)
		}

		if !actionFound {
			eventsToConnector <- NewEvent(event.Type(), "", ErrNoActionFound)
		}
	}
}

func actionThreadPool(eventsToActions <-chan Event, eventsToConnector chan<- Event) {
	for w := 0; w < runtime.NumCPU(); w++ {
		go actionWorker(eventsToActions, eventsToConnector)
	}
}

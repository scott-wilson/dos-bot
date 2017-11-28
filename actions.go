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

func EmitActions(event string, message string, sender User, room Room, connector chan<- Event) {
	connector <- NewEvent(event, message, nil, sender, room)
}

func actionWorker(toActions <-chan Event, toConnector chan<- Event) {
	for event := range toActions {
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
			toConnector <- NewEvent(event.Type(), result, err, event.Sender(), event.Room())
		}

		if !actionFound {
			toConnector <- NewEvent(event.Type(), "", ErrNoActionFound, event.Sender(), event.Room())
		}
	}
}

func actionThreadPool(toActions <-chan Event, toConnector chan<- Event) {
	for w := 0; w < runtime.NumCPU(); w++ {
		go actionWorker(toActions, toConnector)
	}
}

package dosbot

import "log"

func RegisterConnector(connector func(chan<- Event, <-chan Event) func() error) {
	log.Println("Registering connector.")
	connections = append(connections, connector)
}

package dosbot

func RegisterConnector(connector func(Bot, chan<- Event, <-chan Event) func() error) {
	connections = append(connections, connector)
}

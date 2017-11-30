package dosbot

func RegisterConnector(connector func(chan<- Event, <-chan Event) func() error) {
	connections = append(connections, connector)
}

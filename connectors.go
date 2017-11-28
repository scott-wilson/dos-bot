package dosbot

import (
	"bufio"
	"fmt"
	"os"
)

func RegisterConnector(connector func(Bot, chan<- Event, <-chan Event) func() error) {
	connections = append(connections, connector)
}

func TerminalConnector(bot Bot, toActions chan<- Event, toChannel <-chan Event) func() error {
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')

			EmitActions("listen", text, toActions)

			select {
			case result := <-toChannel:
				if err := result.Error(); err != nil {
					if err == ErrNoActionFound {
						continue
					} else {
						fmt.Println("ERROR", err)
						continue
					}
				}

				fmt.Println("RESULT", result.Message())
			}
		}
	}()

	return func() error { return nil }
}

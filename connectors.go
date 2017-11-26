package dosbot

import (
	"bufio"
	"fmt"
	"os"
)

func RegisterConnector(name string, bot Bot, connector func(Bot, chan<- Event, <-chan Event)) {
	input := make(chan Event)
	output := make(chan Event)
	actionThreadPool(input, output)
	connector(bot, input, output)
}

func TerminalConnector(bot Bot, input chan<- Event, output <-chan Event) {
	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')

			EmitActions("listen", text, input)

			select {
			case result := <-output:
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
}

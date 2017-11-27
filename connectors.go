package dosbot

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

func RegisterConnector(name string, bot Bot, connector func(Bot, chan<- Event, <-chan Event) func() error) {
	input := make(chan Event)
	output := make(chan Event)
	actionThreadPool(input, output)

	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	closeFunctions = append(closeFunctions, connector(bot, input, output))
}

func TerminalConnector(bot Bot, input chan<- Event, output <-chan Event) func() error {
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

	return func() error { return nil }
}

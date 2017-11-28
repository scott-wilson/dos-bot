# dos-bot
A simple chat bot in Go(lang) inspired by Hubot

# New Actions
Creating new actions are as simple as this:

```golang
func myAction(input string) (string, error) {
    // Check if action supports input.
    if input != "test" {
        return "", ErrEventNotSupportedByAction
    }

    return "Hello world!", nil
}

func init() {
    RegisterAction("listen", myAction)
}
```
# New Connections
Creating new connections are as simple as this:

```golang
func myConnector(bot Bot, toActions chan<- Event, toChannel <-chan Event) {
    // Inputs
    go func() {
        // Query from service (Discord, Slack, Rocket.Chat, etc) here.
        EmitActions("listen", "test", toActions)
    }()

    // Outputs
    go func() {
        // Output to service from here
    }()
}

func init() {
    RegisterConnector(myConnector)
}
```

# New Bot

```golang
package mybot

import "github.com/scott-wilson/dosbot"

func main() {
    bot := NewBot("test")

    // Register actions first
    dosbot.RegisterAction("listen", myAction)

    // Register connectors
    RegisterConnector(myConnector)

    // Run
    Run(bot)
}
```

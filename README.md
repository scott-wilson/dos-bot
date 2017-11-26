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
func myConnector(bot Bot, input chan<- Event, output <-chan Event) {
    // Inputs
    go func() {
        // Query from service (Discord, Slack, Rocket.Chat, etc) here.
        EmitActions("listen", "test", input)
    }()

    // Outputs
    go func() {
        // Output to service from here
    }()
}

func init() {
    bot := NewBot("test")
    RegisterConnector("myConnector", bot, myConnector)
}
```

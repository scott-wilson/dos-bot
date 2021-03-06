package dosbot

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var connections []func(chan<- Event, <-chan Event) func() error
var closeFunctions []func() error

func Run() {
	log.Println("Starting bot.")

	for _, connector := range connections {
		toActions := make(chan Event)
		toChannel := make(chan Event)
		actionThreadPool(toActions, toChannel)

		mutex := &sync.Mutex{}
		mutex.Lock()
		defer mutex.Unlock()
		closeFunctions = append(closeFunctions, connector(toActions, toChannel))
	}

	defer Close()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func Close() error {
	log.Println("Closing bot.")

	for _, close := range closeFunctions {
		if err := close(); err != nil {
			return err
		}
	}

	return nil
}

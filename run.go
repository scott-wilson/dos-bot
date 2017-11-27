package dosbot

import (
	"os"
	"os/signal"
	"syscall"
)

var closeFunctions []func() error

func Run() {
	defer Close()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func Close() error {
	for _, close := range closeFunctions {
		if err := close(); err != nil {
			return err
		}
	}

	return nil
}

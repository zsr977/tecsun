package helper

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Daemon() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		switch sig {
		case syscall.SIGHUP:
		default:
			time.AfterFunc(time.Duration(int(3e9)), func() {
				os.Exit(1)
			})
			return
		}
	}
}

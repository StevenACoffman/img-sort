package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/StevenACoffman/img-sort/cmd"
	"github.com/StevenACoffman/img-sort/internal/exif"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		exif.Instance().Close()
		os.Exit(1)
	}()

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

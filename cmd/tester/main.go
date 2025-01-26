package main

import (
	"os"
	"os/signal"
	"syscall"

	"wow/internal/pkg/logs"
)

func main() {
	stdin := make(chan os.Signal, 1)
	signal.Notify(stdin, syscall.SIGINT, syscall.SIGTERM)

	<-stdin

	logs.MainLogger.Print("exit")
}

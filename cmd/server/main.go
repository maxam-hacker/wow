package main

import (
	"os"
	"os/signal"
	"syscall"

	"wow/internal/pkg/logs"
	"wow/internal/server"
	"wow/internal/storage"
)

func main() {
	logs.MainLogger.Print("start")

	storage.Initialize("")

	s := server.Server{
		Host: "0.0.0.0",
		Port: 9876,
		WorkLoadBalancer: func(currentWorkLoad int) (int16, error) {
			if currentWorkLoad > 1000000 {
				return 4, nil
			}

			if currentWorkLoad > 10000 {
				return 3, nil
			}

			if currentWorkLoad > 1000 {
				return 2, nil
			}

			return 1, nil
		},
	}

	err := s.Start()
	if err != nil {
		return
	}

	logs.MainLogger.Print("server has started")

	stdin := make(chan os.Signal, 1)
	signal.Notify(stdin, syscall.SIGINT, syscall.SIGTERM)

	<-stdin

	s.Stop()

	logs.MainLogger.Print("exit")
}

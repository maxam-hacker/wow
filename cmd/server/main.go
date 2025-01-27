package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"wow/internal/config"
	"wow/internal/pkg/logs"
	"wow/internal/server"
	"wow/internal/storage"
)

func main() {
	logs.MainServerLogger.Print("start")

	configPath := flag.String("config", "./config/server/config.json", "configuration file path")
	flag.Parse()

	serviceHost := os.Getenv("SERVICE_HOST")
	if serviceHost == "" {
		serviceHost = "0.0.0.0"
	}

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "9877"
	}
	servicePort, err := strconv.Atoi(port)
	if err != nil {
		logs.MainServerLogger.Print("can't get service port", err)
	}

	serverConfig, err := config.NewTcpServerConfiguration(*configPath)
	if err != nil {
		logs.MainServerLogger.Print("can't get configuration", err)
		return
	}

	logs.MainServerLogger.Print("working with configuration", serverConfig)

	storage.Initialize(serverConfig.StorageOpts.PathToBook)

	s := server.Server{
		Opts: server.ServerOpts{
			Host:                serviceHost,
			Port:                servicePort,
			Workres:             10,
			CloseAfterAction:    false,
			CloseAfterExecution: false,
		},
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

	err = s.Start()
	if err != nil {
		return
	}

	logs.MainServerLogger.Print("server has started")

	stdin := make(chan os.Signal, 1)
	signal.Notify(stdin, syscall.SIGINT, syscall.SIGTERM)

	<-stdin

	s.Stop()

	logs.MainServerLogger.Print("exit")
}

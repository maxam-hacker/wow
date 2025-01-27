package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {
	go TestManyConnections()

	stdin := make(chan os.Signal, 1)
	signal.Notify(stdin, syscall.SIGINT, syscall.SIGTERM)

	<-stdin
}

func TestManyConnections() {
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
		fmt.Print("can't get service port", err)
	}

	for idx := 0; idx < 20000; idx++ {
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serviceHost, servicePort))
			if err != nil {
				return
			}

			var rxBuff [1024]byte

			_, err = conn.Read(rxBuff[:])
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

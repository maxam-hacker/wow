package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"wow/internal/client"
	"wow/internal/pkg/logs"
	"wow/internal/proto"
	"wow/pkg/logger"
)

var (
	ErrUnknownResponse = errors.New("unknown response")
)

type Client struct {
	Id          string
	Version     string
	MaxLines    int
	lineCounter int
	tcpClient   client.TcpClient
}

func main() {
	c := &Client{
		Id:          uuid.NewString(),
		Version:     "0.0.1",
		MaxLines:    20,
		lineCounter: 1,
		tcpClient: client.TcpClient{
			Host: "127.0.0.1",
		},
	}

	c.tcpClient.SetRxHandler(c.RxHandler)

	go c.tcpClient.Start()

	go c.Loop()

	logs.MainClientLogger.Print("client has started")

	stdin := make(chan os.Signal, 1)
	signal.Notify(stdin, syscall.SIGINT, syscall.SIGTERM)

	<-stdin

	c.tcpClient.Stop()
}

func (c *Client) Loop() {
	for {
		c.tcpClient.SendRequestOnAction(c.Id, c.Version)
		time.Sleep(10 * time.Second)
	}
}

func (c *Client) RxHandler(responseBytes []byte) error {
	var response proto.Response

	err := json.Unmarshal(responseBytes, &response)
	if err != nil {
		return err
	}

	if response.Type == proto.ResponseOnActionType {
		logs.MainClientLogger.PrintWithContext("action", logger.Context{
			"response:": response,
		})

		hash, err := response.Hash.Compute(10000000)
		if err != nil {
			logs.MainClientLogger.PrintWithContext("can't compute hashcash", logger.Context{
				"error":     err,
				"response:": response,
			})
			return err
		}

		c.tcpClient.SendRequestOnActionExecution(c.Id, c.Version, hash, c.lineCounter)

	} else if response.Type == proto.ResponseOnActionExecutionType {
		c.lineCounter++

		c.lineCounter %= c.MaxLines

		logs.MainClientLogger.PrintWithContext("execution", logger.Context{
			"LineId:":    response.LineId,
			"LineValue:": response.Result,
		})

	} else {
		err = ErrUnknownResponse
		logs.MainClientLogger.PrintWithContext("can't handle response", logger.Context{
			"error":     err,
			"response:": response,
		})
		return err
	}

	return nil
}

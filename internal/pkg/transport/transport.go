package transport

import (
	"errors"

	"wow/internal/pkg/transport/epoll"
	"wow/internal/pkg/transport/workers"
	"wow/internal/types"
)

const (
	DefaultWorkersNumber = 10
)

var (
	ErrUnimplementedFunctionality = errors.New("unimplemented")
	ErrEmptyMessageHandler        = errors.New("empty message handler")
)

type TcpServer struct {
	Host           string
	Port           int
	WorkersNumber  int
	MessageHandler types.TcpServerMessageHandler
	epoll          *epoll.Epoll
	workers        *workers.Pool
}

func (tcpServer *TcpServer) Start(epollOpts epoll.EpollOpts) error {
	var err error

	if tcpServer.WorkersNumber == 0 {
		tcpServer.WorkersNumber = DefaultWorkersNumber
	}

	if tcpServer.MessageHandler == nil {
		return ErrEmptyMessageHandler
	}

	tcpServer.workers, err = workers.New(tcpServer.WorkersNumber, tcpServer.MessageHandler)
	if err != nil {
		return err
	}

	tcpServer.epoll, err = epoll.New(tcpServer.Host, tcpServer.Port, tcpServer.workers, epollOpts)
	if err != nil {
		return err
	}

	tcpServer.workers.Start(tcpServer.epoll.Delete)

	tcpServer.epoll.Start()

	return nil
}

func (tcpServer *TcpServer) Stop() {
	tcpServer.epoll.Stop()
	tcpServer.workers.Stop()
}

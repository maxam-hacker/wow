package server

import (
	"encoding/json"
	"errors"
	"io"

	"wow/internal/pkg/logs"
	"wow/internal/pkg/transport"
	"wow/internal/pkg/transport/epoll"
	"wow/internal/proto"
	"wow/internal/storage"
	"wow/pkg/logger"
)

type ServerOpts struct {
	Host                string `json:"host"`
	Port                int    `json:"port"`
	Workres             int    `json:"workers"`
	CloseAfterAction    bool   `json:"closeAfterAction"`
	CloseAfterExecution bool   `json:"closeAfterExecution"`
}

type Server struct {
	WorkLoadBalancer func(int) (int16, error)
	tcpServer        transport.TcpServer
	Opts             ServerOpts
	EpollOpts        epoll.EpollOpts
}

var (
	ErrEmptyLoadBalancer = errors.New("can't find load balancer algo")
	ErrEmptyWriter       = errors.New("can't find writer")
	ErrEmptyCloser       = errors.New("can't find closer")
	ErrUnknownRequest    = errors.New("unknown request")
	ErrValiadionError    = errors.New("validation error")
)

func (server *Server) Start() error {
	if server.WorkLoadBalancer == nil {
		return ErrEmptyLoadBalancer
	}

	server.tcpServer = transport.TcpServer{
		Host:           server.Opts.Host,
		Port:           server.Opts.Port,
		MessageHandler: server.messageHandler,
	}

	return server.tcpServer.Start(server.EpollOpts)
}

func (server *Server) Stop() {
	server.tcpServer.Stop()
}

func (server *Server) messageHandler(requestBytes []byte, writer io.Writer, currentWorkLoad int, closer io.Closer) error {
	var request proto.Request
	var err error

	errContext := logger.Context{
		"requestBytes":    requestBytes,
		"currentWorkLoad": currentWorkLoad,
		"requestAsString": string(requestBytes),
	}

	if writer == nil {
		return ErrEmptyWriter
	}

	if closer == nil {
		return ErrEmptyCloser
	}

	err = json.Unmarshal(requestBytes, &request)
	if err != nil {
		errContext["error"] = err
		logs.ServerLogger.PrintWithContext("unmarshalling error", errContext)
		closer.Close()
		return err
	}

	if request.Type == proto.RequestActionType {
		workLoadFactor, err := server.getWorkLoadFactor(currentWorkLoad)
		if err != nil {
			errContext["error"] = err
			logs.ServerLogger.PrintWithContext("can't calculate workload factor", errContext)
			closer.Close()
			return err
		}

		response := proto.NewResponseOnAction(workLoadFactor)

		responseBytes, err := json.Marshal(response)
		if err != nil {
			errContext["error"] = err
			errContext["response"] = response
			logs.ServerLogger.PrintWithContext("marshalling error", errContext)
			closer.Close()
			return err
		}

		_, err = writer.Write(responseBytes)
		if err != nil {
			errContext["error"] = err
			errContext["responseBytes"] = responseBytes
			logs.ServerLogger.PrintWithContext("can't write data", errContext)
			closer.Close()
			return err
		}

		if server.Opts.CloseAfterAction {
			closer.Close()
		}

	} else if request.Type == proto.RequestActionExecutionType {
		workLoadFactor, err := server.getWorkLoadFactor(currentWorkLoad)
		if err != nil {
			errContext["error"] = err
			logs.ServerLogger.PrintWithContext("can't calculate workload factor", errContext)
			closer.Close()
			return err
		}

		ok, err := request.Validate(workLoadFactor)
		if err != nil {
			errContext["error"] = err
			errContext["workLoadFactor"] = workLoadFactor
			logs.ServerLogger.PrintWithContext("can't calculate workload factor", errContext)
			closer.Close()
			return err
		}

		if !ok {
			errContext["workLoadFactor"] = workLoadFactor
			logs.ServerLogger.PrintWithContext("can't calculate workload factor", errContext)
			closer.Close()
			return ErrValiadionError
		}

		response := proto.NewResponseOnActionExecution(request.LineId, storage.GetLine(request.LineId))

		responseBytes, err := json.Marshal(response)
		if err != nil {
			errContext["error"] = err
			errContext["response"] = response
			logs.ServerLogger.PrintWithContext("marshalling error", errContext)
			closer.Close()
			return err
		}

		_, err = writer.Write(responseBytes)
		if err != nil {
			errContext["error"] = err
			errContext["responseBytes"] = responseBytes
			logs.ServerLogger.PrintWithContext("can't write data", errContext)
			closer.Close()
			return err
		}

		if server.Opts.CloseAfterExecution {
			closer.Close()
		}

	} else {
		logs.ServerLogger.PrintWithContext("unknown request", errContext)
		closer.Close()
		return ErrUnknownRequest
	}

	return nil
}

func (server *Server) getWorkLoadFactor(currentWorkLoad int) (int16, error) {
	workLoadFactor, err := server.WorkLoadBalancer(currentWorkLoad)
	if err != nil {
		logs.ServerLogger.Print("can't get load factor", err)
		return 1000, err
	}

	return workLoadFactor, nil
}

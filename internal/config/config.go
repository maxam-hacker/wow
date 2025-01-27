package config

import (
	"wow/internal/pkg/transport/epoll"
	"wow/internal/server"
)

type TcpServerConfiguration struct {
	ServerOpts server.ServerOpts `json:"server"`
	EpollOpts  epoll.EpollOpts   `json:"epoll"`
}

func NewTcpServerConfiguration(configPath string) *TcpServerConfiguration {
	return &TcpServerConfiguration{}
}

type TcpClientConfiguration struct {
}

func NewTcpClientConfiguration(configPath string) *TcpClientConfiguration {
	return &TcpClientConfiguration{}
}

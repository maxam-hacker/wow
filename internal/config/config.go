package config

import (
	"wow/internal/pkg/transport"
	"wow/internal/pkg/transport/epoll"
)

type TcpServerConfiguration struct {
	TcpServerOpts transport.TcpServerOpts `json:"server"`
	EpollOpts     epoll.EpollOpts         `json:"epoll"`
}

func NewTcpServerConfiguration(configPath string) *TcpServerConfiguration {
	return &TcpServerConfiguration{}
}

type TcpClientConfiguration struct {
}

func NewTcpClientConfiguration(configPath string) *TcpClientConfiguration {
	return &TcpClientConfiguration{}
}

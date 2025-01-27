package config

import (
	"wow/internal/pkg/transport/epoll"
	"wow/internal/server"
	"wow/internal/storage"
	"wow/internal/types"
)

type TcpServerConfiguration struct {
	ServerOpts           server.ServerOpts          `json:"server"`
	EpollOpts            epoll.EpollOpts            `json:"epoll"`
	StorageOpts          storage.StorageOpts        `json:"storage"`
	WorkLoadBalancerOpts types.WorkLoadBalancerOpts `json:"workLoadBalancer"`
}

func NewTcpServerConfiguration(configPath string) *TcpServerConfiguration {
	return &TcpServerConfiguration{}
}

package config

import (
	"encoding/json"
	"os"
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

func NewServerConfiguration(configPath string) (*TcpServerConfiguration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var serverConfig TcpServerConfiguration

	err = json.Unmarshal(data, &serverConfig)
	if err != nil {
		return nil, err
	}

	return &serverConfig, nil
}

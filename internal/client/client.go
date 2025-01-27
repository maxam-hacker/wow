package client

import (
	"encoding/json"
	"fmt"
	"net"

	"wow/internal/hashcash"
	"wow/internal/proto"
	"wow/internal/types"
)

type TcpClient struct {
	Host       string
	Port       int
	Connection net.Conn
	RxHandler  types.TcpClientMessageHandler
}

var RxBuffer [65536]byte

func (client *TcpClient) Start() error {
	var err error

	client.Connection, err = net.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer client.Connection.Close()

	for {
		n, _ := client.Connection.Read(RxBuffer[:])
		if n > 0 {
			if client.RxHandler != nil {
				client.RxHandler(RxBuffer[:n])
			}
		}
	}

}

func (client *TcpClient) SetRxHandler(rxHandler types.TcpClientMessageHandler) {
	client.RxHandler = rxHandler
}

func (client *TcpClient) SendRequestOnAction(clientId string, clientVersion string) (int, error) {
	if client.Connection == nil {
		return 0, nil
	}

	request := proto.NewRequestAction(
		proto.ClientMeta{
			ClientId:      clientId,
			ClientVersion: clientVersion,
		},
	)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	return client.Connection.Write(requestBytes)

}

func (client *TcpClient) SendRequestOnActionExecution(clientId string, clientVersion string, hash hashcash.Hashcash, lineId int) (int, error) {
	if client.Connection == nil {
		return 0, nil
	}

	request := proto.NewRequestActionExecution(
		proto.ClientMeta{
			ClientId:      clientId,
			ClientVersion: clientVersion,
		},
		lineId,
		hash,
	)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	return client.Connection.Write(requestBytes)

}

func (client *TcpClient) Stop() {
	client.Connection.Close()
}

package types

import "io"

type WorkLoadBalancerOpts struct {
	MinZeros              int `json:"minZeros"`
	Threshold1Connections int `json:"threshold1Connections"`
	Threshold1Zeros       int `json:"threshold1Zeros"`
	Threshold2Connections int `json:"threshold2Connections"`
	Threshold2Zeros       int `json:"threshold2Zeros"`
	Threshold3Connections int `json:"threshold3Connections"`
	Threshold3Zeros       int `json:"threshold3Zeros"`
}

type TcpServerMessageHandler func([]byte, io.Writer, int, io.Closer) error

type TcpClientMessageHandler func([]byte) error

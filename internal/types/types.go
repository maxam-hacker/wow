package types

import "io"

type TcpServerMessageHandler func([]byte, io.Writer, int) error

type TcpClientMessageHandler func([]byte) error

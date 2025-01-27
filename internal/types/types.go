package types

import "io"

type TcpServerMessageHandler func([]byte, io.Writer, int, io.Closer) error

type TcpClientMessageHandler func([]byte) error

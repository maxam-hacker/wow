package workers

import (
	"golang.org/x/sys/unix"

	"wow/internal/types"
)

type Worker struct {
	Id        int
	InputGate chan PoolMessage
	Buffer    []byte
	Handler   types.TcpServerMessageHandler
	Writer    Writer
	Closer    Closer
	Cancel    chan struct{}
}

func (worker *Worker) start() {
	var err error
	var n int

	for {
		select {
		case poolMessage := <-worker.InputGate:
			n, _ = unix.Read(poolMessage.TargetSocketHandler, worker.Buffer)
			if n > 0 {
				if worker.Handler != nil {
					worker.Writer.SetTarget(poolMessage.TargetSocketHandler)
					worker.Closer.SetTarget(poolMessage.TargetSocketHandler)
					err = worker.Handler(worker.Buffer[0:n], worker.Writer, poolMessage.ConnectionsNum, worker.Closer)
					if err != nil {
						unix.Close(poolMessage.TargetSocketHandler)
					}
				}
			}
		case <-worker.Cancel:
			return
		}
	}
}

func (worker *Worker) stop() {
	worker.Cancel <- struct{}{}
}

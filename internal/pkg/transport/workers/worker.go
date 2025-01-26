package workers

import (
	"golang.org/x/sys/unix"

	"wow/internal/pkg/logs"
	"wow/internal/types"
)

type Worker struct {
	Id        int
	InputGate chan PoolMessage
	Buffer    []byte
	Handler   types.TcpServerMessageHandler
	Writer    Writer
}

func (worker *Worker) start() {
	logs.WorkersWorkerLogger.Print("start worker", worker.Id)
	defer logs.WorkersWorkerLogger.Print("start worker : done", worker.Id)

	var err error
	var n int

	for poolMessage := range worker.InputGate {
		n, _ = unix.Read(poolMessage.TargetSocketHandler, worker.Buffer)
		if n > 0 {
			if worker.Handler != nil {
				worker.Writer.SetTarget(poolMessage.TargetSocketHandler)

				err = worker.Handler(worker.Buffer[0:n], worker.Writer, poolMessage.ConnectionsNum)
				if err != nil {
					unix.Close(poolMessage.TargetSocketHandler)
				}
			}
		}
	}
}

func (worker *Worker) stop() {
	logs.WorkersWorkerLogger.Print("stop worker", worker.Id)
	logs.WorkersWorkerLogger.Print("stop worker : done", worker.Id)
}

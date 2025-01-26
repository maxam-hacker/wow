package workers

import (
	"wow/internal/pkg/logs"
	"wow/internal/types"
)

type PoolMessage struct {
	TargetSocketHandler int
	ConnectionsNum      int
}

type Pool struct {
	Workres        []*Worker
	WorkersNumber  int
	InputGate      chan PoolMessage
	messageHandler types.TcpServerMessageHandler
}

func New(workersNumber int, messageHandler types.TcpServerMessageHandler) (*Pool, error) {
	logs.WorkersLogger.Print("new workers pool", workersNumber, messageHandler)

	wp := &Pool{
		WorkersNumber:  workersNumber,
		InputGate:      make(chan PoolMessage),
		messageHandler: messageHandler,
	}

	err := wp.startWorkers()
	if err != nil {
		return nil, err
	}

	return wp, nil
}

func (pool *Pool) HandleMessage(targetSocketHandler int, connectionsNum int) {
	pool.InputGate <- PoolMessage{
		TargetSocketHandler: targetSocketHandler,
		ConnectionsNum:      connectionsNum,
	}
}

func (pool *Pool) startWorkers() error {
	logs.WorkersPoolLogger.Print("start workers", pool)

	for idx := range pool.WorkersNumber {
		w := &Worker{
			Id:        idx,
			InputGate: pool.InputGate,
			Buffer:    make([]byte, 65536),
			Handler:   pool.messageHandler,
		}

		pool.Workres = append(pool.Workres, w)

		go w.start()
	}

	logs.WorkersPoolLogger.Print("start workers : done", pool)

	return nil
}

func (pool *Pool) Stop() {
	for _, worker := range pool.Workres {
		worker.stop()
	}
}

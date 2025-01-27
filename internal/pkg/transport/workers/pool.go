package workers

import (
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
	epCloser       func(int) error
}

func New(workersNumber int, messageHandler types.TcpServerMessageHandler) (*Pool, error) {
	wp := &Pool{
		WorkersNumber:  workersNumber,
		InputGate:      make(chan PoolMessage),
		messageHandler: messageHandler,
	}

	return wp, nil
}

func (pool *Pool) Start(epCloser func(int) error) error {
	pool.epCloser = epCloser

	err := pool.startWorkers()
	if err != nil {
		return err
	}

	return nil
}

func (pool *Pool) HandleMessage(targetSocketHandler int, connectionsNum int) {
	pool.InputGate <- PoolMessage{
		TargetSocketHandler: targetSocketHandler,
		ConnectionsNum:      connectionsNum,
	}
}

func (pool *Pool) startWorkers() error {
	for idx := range pool.WorkersNumber {
		w := &Worker{
			Id:        idx,
			InputGate: pool.InputGate,
			Buffer:    make([]byte, 65536),
			Handler:   pool.messageHandler,
		}

		w.Closer.epCloser = pool.epCloser

		pool.Workres = append(pool.Workres, w)

		go w.start()
	}

	return nil
}

func (pool *Pool) Stop() {
	for _, worker := range pool.Workres {
		worker.stop()
	}
}

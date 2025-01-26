package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	acceptedConnectionsCounterMu sync.RWMutex
	acceptedConnectionsCounter   int

	readConnectionsCounterMu sync.RWMutex
	readConnectionsCounter   int

	Connections sync.Map
}

func (metrics *Metrics) IncConnections(targetSocketHandler int) {
	metrics.acceptedConnectionsCounterMu.Lock()
	defer metrics.acceptedConnectionsCounterMu.Unlock()

	connectionElement := &ConnectionElement{
		SocketHandler: targetSocketHandler,
		LastActivity:  time.Now().UTC(),
	}

	metrics.acceptedConnectionsCounter++

	metrics.Connections.Store(targetSocketHandler, connectionElement)
}

func (metrics *Metrics) DecConnections(targetSocketHandler int) {
	metrics.acceptedConnectionsCounterMu.Lock()
	defer metrics.acceptedConnectionsCounterMu.Unlock()

	metrics.acceptedConnectionsCounter--

	_, exists := metrics.Connections.Load(targetSocketHandler)
	if exists {
		metrics.Connections.Delete(targetSocketHandler)
	}
}

func (metrics *Metrics) IncRead(targetSocketHandler int) {
	metrics.readConnectionsCounterMu.Lock()
	defer metrics.readConnectionsCounterMu.Unlock()

	metrics.readConnectionsCounter++

	connectionElement, exists := metrics.Connections.Load(targetSocketHandler)
	if exists {
		connectionElement.(*ConnectionElement).LastActivity = time.Now().UTC()
	}
}

func (metrics *Metrics) Walk(processor func(connection *ConnectionElement)) {
	metrics.Connections.Range(func(k, v any) bool {
		processor(v.(*ConnectionElement))
		return true
	})
}

func (metrics *Metrics) GetConnections() int {
	metrics.acceptedConnectionsCounterMu.RLock()
	defer metrics.acceptedConnectionsCounterMu.RUnlock()

	return metrics.acceptedConnectionsCounter
}

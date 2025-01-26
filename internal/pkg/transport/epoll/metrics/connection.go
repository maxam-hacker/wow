package metrics

import "time"

type ConnectionElement struct {
	SocketHandler int
	LastActivity  time.Time
}

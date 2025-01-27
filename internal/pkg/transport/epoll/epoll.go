package epoll

import (
	"net"
	"time"

	"golang.org/x/sys/unix"

	"wow/internal/pkg/logs"
	"wow/internal/pkg/transport/epoll/metrics"
	"wow/internal/pkg/transport/workers"
	"wow/pkg/logger"
)

const (
	DefaultEpollEventsBufferSize       = 128
	DefaultEpollLoopWaitTimeout        = 5000
	DefaultCleanerPeriod               = 30
	DefaultCleanerConnectionsThreshold = 100000
	DefaultCleanerTimeThreshold        = 180
)

type EpollOpts struct {
	EpollEventsBufferSize       int
	EpollLoopWaitTimeout        int
	CleanerPeriod               int
	CleanerConnectionsThreshold int
	CleanerTimeThreshold        int
}

type Epoll struct {
	Host                  string
	Port                  int
	EpollEventsBufferSize int
	EpollLoopWaitTimeout  int
	epollEventsBuffer     []unix.EpollEvent
	epollHandler          int
	socketHandler         int
	workers               *workers.Pool
	metrics               metrics.Metrics
}

func New(host string, port int, workers *workers.Pool) (*Epoll, error) {
	return &Epoll{
		Host:                  host,
		Port:                  port,
		EpollEventsBufferSize: DefaultEpollEventsBufferSize,
		EpollLoopWaitTimeout:  DefaultEpollLoopWaitTimeout,
		workers:               workers,
		metrics:               metrics.Metrics{},
	}, nil
}

func (ep *Epoll) Start() error {
	err := ep.initialize()
	if err != nil {
		return err
	}

	go ep.loop()
	go ep.cleaner()

	return nil
}

func (ep *Epoll) Stop() {
}

func (ep *Epoll) initialize() error {
	var err error

	erContext := logger.Context{
		"epoll": ep,
	}

	ep.epollHandler, err = unix.EpollCreate1(0)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't create epoll", erContext)
		return err
	}

	ep.epollEventsBuffer = make([]unix.EpollEvent, ep.EpollEventsBufferSize)

	ep.socketHandler, err = unix.Socket(unix.AF_INET, unix.O_NONBLOCK|unix.SOCK_STREAM, 0)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't create socket", erContext)
		return err
	}

	err = unix.SetNonblock(ep.socketHandler, true)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't set nonblocking socket optinon", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.SetsockoptInt(ep.socketHandler, unix.IPPROTO_TCP, unix.TCP_NODELAY, 1)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't set no delay socket optinon", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.SetsockoptInt(ep.socketHandler, unix.IPPROTO_TCP, unix.TCP_QUICKACK, 1)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't set quick ack socket optinon", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.SetsockoptInt(ep.socketHandler, unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't set reuse addr socket optinon", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.SetsockoptInt(ep.socketHandler, unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't set reuse port socket optinon", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	addr := unix.SockaddrInet4{
		Port: ep.Port,
	}
	copy(addr.Addr[:], net.ParseIP(ep.Host).To4())

	if err = unix.Bind(ep.socketHandler, &addr); err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't bind socket", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.Listen(ep.socketHandler, unix.SOMAXCONN)
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't listen socket", erContext)
		unix.Close(ep.socketHandler)
		return err
	}

	err = unix.EpollCtl(ep.epollHandler, unix.EPOLL_CTL_ADD, ep.socketHandler, &unix.EpollEvent{
		Fd:     int32(ep.socketHandler),
		Events: unix.EPOLLIN | unix.EPOLLHUP | unix.EPOLLRDHUP | unix.EPOLLET | unix.EPOLLERR,
	})
	if err != nil {
		erContext["error"] = err
		logs.EpollLogger.PrintWithContext("can't add socket in epoll", erContext)
		return err
	}

	return nil
}

func (ep *Epoll) loop() {
	logs.EpollLogger.Print("epoll control loop", ep)

	for {
		n, err := unix.EpollWait(ep.epollHandler, ep.epollEventsBuffer, ep.EpollLoopWaitTimeout)
		if err == nil {
			if n == 0 {
				continue
			}

			for idx := 0; idx < n; idx++ {
				fd := int(ep.epollEventsBuffer[idx].Fd)
				events := ep.epollEventsBuffer[idx].Events

				if fd == ep.socketHandler {
					ep.accept()
					continue
				}

				if events&(unix.EPOLLHUP|unix.EPOLLRDHUP) != 0 {
					ep.Delete(fd)
					continue
				}

				if events&unix.EPOLLERR != 0 {
					ep.Delete(fd)
					continue
				}

				if events&unix.EPOLLIN != 0 {
					ep.read(fd)
				} else {
					logs.EpollLogger.Print("loop: unexpectable behaviour", fd)
					ep.Delete(fd)
				}
			}
		} else {
			if err != unix.EINTR {
				break
			}
		}
	}
}

func (ep *Epoll) accept() error {
	erContext := logger.Context{}

	for {
		acceptedSocketHandler, _, err := unix.Accept(ep.socketHandler)
		if err != nil {
			return err
		}

		erContext["epoll"] = ep
		erContext["acceptedSocketHandler"] = acceptedSocketHandler

		ep.metrics.IncConnections(acceptedSocketHandler)

		err = unix.SetNonblock(acceptedSocketHandler, true)
		if err != nil {
			erContext["error"] = err
			logs.EpollLogger.PrintWithContext("can't set nonblocking option to accepted socket", erContext)
			return err
		}

		err = unix.SetsockoptInt(acceptedSocketHandler, unix.SOL_SOCKET, unix.SO_KEEPALIVE, 0)
		if err != nil {
			erContext["error"] = err
			logs.EpollLogger.PrintWithContext("can't set keep alive option to accepted socket", erContext)
			return err
		}

		err = unix.EpollCtl(ep.epollHandler, unix.EPOLL_CTL_ADD, acceptedSocketHandler, &unix.EpollEvent{
			Events: unix.EPOLLIN | unix.EPOLLHUP | unix.EPOLLRDHUP | unix.EPOLLET | unix.EPOLLERR,
			Fd:     int32(acceptedSocketHandler),
		})
		if err != nil {
			erContext["error"] = err
			logs.EpollLogger.PrintWithContext("can't add accepted socket in epoll", erContext)
			return err
		}

	}
}

func (ep *Epoll) read(targetSocketHandler int) error {
	ep.metrics.IncRead(targetSocketHandler)

	connectionsNum := ep.metrics.GetConnections()

	ep.workers.HandleMessage(targetSocketHandler, connectionsNum)

	return nil
}

func (ep *Epoll) Delete(targetSocketHandler int) error {
	err := unix.EpollCtl(ep.epollHandler, unix.EPOLL_CTL_DEL, targetSocketHandler, nil)
	if err != nil {
		return err
	}

	unix.Close(targetSocketHandler)

	ep.metrics.DecConnections(targetSocketHandler)

	return nil
}

func (ep *Epoll) cleaner() {
	for {
		time.Sleep(DefaultCleanerPeriod * time.Second)

		connectionsNum := ep.metrics.GetConnections()

		logs.EpollLogger.PrintWithContext("cleaner: metrics", logger.Context{
			"Connections": ep.metrics.GetConnections(),
		})

		if connectionsNum < DefaultCleanerConnectionsThreshold {
			continue
		}

		ep.metrics.Walk(func(connectionElement *metrics.ConnectionElement) {
			if time.Now().UTC().Sub(connectionElement.LastActivity) > DefaultCleanerTimeThreshold*time.Second {
				ep.Delete(connectionElement.SocketHandler)
			}
		})
	}
}

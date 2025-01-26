package workers

import "golang.org/x/sys/unix"

type Writer struct {
	targetSocketHandler int
}

func (w *Writer) SetTarget(targetSocketHandler int) {
	w.targetSocketHandler = targetSocketHandler
}

func (w Writer) Write(data []byte) (int, error) {
	return unix.Write(w.targetSocketHandler, data)
}

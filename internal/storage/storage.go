package storage

import "sync"

var storage sync.Map

func Initialize(pathToBook string) {
	storage.Store(1, "line 1")
	storage.Store(2, "line 2")
	storage.Store(3, "line 3")
	storage.Store(4, "line 4")
	storage.Store(5, "line 5")
	storage.Store(6, "line 6")
	storage.Store(7, "line 7")
	storage.Store(8, "line 8")
	storage.Store(9, "line 9")
}

func GetLine(lineId int) string {
	value, ok := storage.Load(lineId)
	if !ok {
		return ""
	}

	return value.(string)
}

package storage

import (
	"bufio"
	"os"
	"sync"
)

type StorageOpts struct {
	PathToBook string `json:"pathToBook"`
}

var storage sync.Map

func Initialize(pathToBook string) error {
	file, err := os.Open(pathToBook)
	if err != nil {
		return err
	}
	defer file.Close()

	idx := 1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		storage.Store(idx, scanner.Text())
		idx++
	}

	return nil
}

func GetLine(lineId int) string {
	value, ok := storage.Load(lineId)
	if !ok {
		return ""
	}

	return value.(string)
}

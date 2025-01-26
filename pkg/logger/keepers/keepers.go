package keepers

import (
	"fmt"
	"os"
	"path"
	"sync"
)

type Message struct {
	keeper  *Keeper
	Prefix1 string
	Prefix2 string
	Prefix3 string
	Text    string
	At      string
}

var serviceLogFile *os.File
var serviceQueue = make(chan Message, 256)
var serviceWaiter sync.WaitGroup

func Initialize(basePath string) {
	if serviceLogFile != nil {
		return
	}

	baseDirectoryPath := path.Join(basePath)

	err := os.MkdirAll(baseDirectoryPath, os.ModePerm)
	if err != nil {
		fmt.Println("can't create directory for service log file", basePath)
		return
	} else {
		serviceFilePath := path.Join(basePath, "service-part000000.log")

		serviceLogFile, err = os.Create(serviceFilePath)
		if err != nil {
			fmt.Println("can't create service log file", serviceLogFile)
			return
		}
	}

	go Start()
}

func New(basePath string) *Keeper {
	Initialize(basePath)

	return &Keeper{
		basePath:  basePath,
		fileCache: make(map[string]*os.File),
	}
}

func Start() {
	serviceWaiter.Add(1)
	defer serviceWaiter.Done()

	for message := range serviceQueue {
		SaveToServiceFile(&message)
		if message.keeper != nil {
			message.keeper.SaveToModuleFile(&message)
		}
	}
}

func SaveToServiceFile(message *Message) {
	if serviceLogFile == nil {
		return
	}

	prefixes := path.Join(message.Prefix1, message.Prefix2, message.Prefix3)

	_, err := serviceLogFile.WriteString(message.At + " ::: " + prefixes + " ::: " + message.Text + "\n")
	if err != nil {
		fmt.Println("can't write to service log file")
	}
}

func Close() {
	close(serviceQueue)
	serviceWaiter.Wait()

	if serviceLogFile != nil {
		serviceLogFile.Close()
	}
}

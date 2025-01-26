package keepers

import (
	"fmt"
	"os"
	"path"
)

type Keeper struct {
	basePath  string
	fileCache map[string]*os.File
}

func (keeper *Keeper) Ingest(message Message) {
	message.keeper = keeper
	serviceQueue <- message
}

func (keeper *Keeper) SaveToModuleFile(message *Message) {
	targetFileDirectoryPath := path.Join(keeper.basePath, message.Prefix1, message.Prefix2, message.Prefix3)
	targetFilePath := path.Join(targetFileDirectoryPath, "module-part000000.log")

	var targetFile *os.File

	cachedFile, exists := keeper.fileCache[targetFilePath]
	if !exists {
		err := os.MkdirAll(targetFileDirectoryPath, os.ModePerm)
		if err != nil {
			fmt.Println("can't create directory for log file", targetFilePath)
			return
		}

		newFile, err := os.Create(targetFilePath)
		if err != nil {
			fmt.Println("can't create log file", targetFilePath)
			return
		}

		keeper.fileCache[targetFilePath] = newFile

		targetFile = newFile

	} else {
		targetFile = cachedFile
	}

	_, err := targetFile.WriteString(message.At + " ::: " + message.Text + "\n")
	if err != nil {
		fmt.Println("can't write to log file", targetFilePath)
	}
}

func (keeper *Keeper) Close() {
	for _, file := range keeper.fileCache {
		file.Close()
	}
}

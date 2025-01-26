package logger

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/google/uuid"

	"wow/pkg/logger/keepers"
)

type Logger struct {
	id         string
	baseLogger log.Logger
	prefix1    string
	prefix2    string
	prefix3    string
	keeper     *keepers.Keeper
}

type Context map[string]interface{}

var loggersCache map[string]*Logger = make(map[string]*Logger)

func New() *Logger {
	l := &Logger{
		id:         uuid.New().String(),
		baseLogger: *log.New(os.Stdout, "", log.Ldate|log.Ltime),
	}

	loggersCache[l.id] = l

	return l
}

func NewWithKeeping(basePath string) *Logger {
	l := &Logger{
		id:         uuid.New().String(),
		baseLogger: *log.New(os.Stdout, "", log.Ldate|log.Ltime),
		keeper:     keepers.New(basePath),
	}

	loggersCache[l.id] = l

	return l
}

func Close() {
	keepers.Close()

	for _, logger := range loggersCache {
		logger.close()
	}
}

func (logger *Logger) WithPrefix(prefix string) *Logger {
	logger.prefix1 = prefix
	logger.baseLogger.SetPrefix(logger.prefix1 + " :: ")
	return logger
}

func (logger *Logger) WithSecondPrefix(secondPrefix string) *Logger {
	logger.prefix2 = secondPrefix
	logger.baseLogger.SetPrefix(logger.prefix1 + " :: " + logger.prefix2 + " :: ")
	return logger
}

func (logger *Logger) WithThirdPrefix(thirdPrefix string) *Logger {
	logger.prefix3 = thirdPrefix
	logger.baseLogger.SetPrefix(logger.prefix1 + " :: " + logger.prefix2 + " :: " + logger.prefix3 + " :: ")
	return logger
}

func (logger *Logger) Print(message string, simpleContext ...any) {
	contextText := ""
	callerText := ""

	for idx, oneContextPart := range simpleContext {
		partTypeName := reflect.TypeOf(oneContextPart).String()
		contextText += fmt.Sprintf("\tcontext item (%d): %s, %+v;\n", idx, partTypeName, oneContextPart)
	}

	if contextText != "" {
		contextText = "\n" + contextText

		_, file, line, ok := runtime.Caller(1)
		if ok {
			callerText += fmt.Sprintf("\t\twhere: %s, line: %d;", file, line)
		}
	}

	text := fmt.Sprintf("%s;%s%s",
		message,
		contextText,
		callerText,
	)

	logger.baseLogger.Println(text)

	if logger.keeper != nil {
		logger.keeper.Ingest(keepers.Message{
			Prefix1: logger.prefix1,
			Prefix2: logger.prefix2,
			Prefix3: logger.prefix3,
			Text:    text,
			At:      time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func (logger *Logger) PrintWithContext(message string, context Context) {
	contextText := ""
	callerText := ""

	for oneContextName, oneContextPart := range context {
		partTypeName := reflect.TypeOf(oneContextPart).String()
		contextText += fmt.Sprintf("\tcontext item [%s]: %s, %+v;\n", oneContextName, partTypeName, oneContextPart)
	}

	_, file, line, ok := runtime.Caller(1)
	if ok {
		callerText += fmt.Sprintf("\t\twhere: %s, line: %d;", file, line)
	}

	text := fmt.Sprintf("%s;\n%s%s",
		message,
		contextText,
		callerText,
	)

	logger.baseLogger.Println(text)

	if logger.keeper != nil {
		logger.keeper.Ingest(keepers.Message{
			Prefix1: logger.prefix1,
			Prefix2: logger.prefix2,
			Prefix3: logger.prefix3,
			Text:    text,
			At:      time.Now().UTC().Format(time.RFC3339),
		})
	}
}

func (logger *Logger) close() {
	if logger.keeper != nil {
		logger.keeper.Close()
	}
}

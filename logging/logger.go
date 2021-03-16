package logging

import (
	"io"
	"log"
	"sync"
)

type Source struct {
	logLevel string
	mutex    sync.Mutex
}

func CreateSource(logLevel string) *Source {
	source := Source{logLevel: logLevel}
	return &source
}

func (s *Source) CreateLogger(out io.Writer, prefix string, flag int) *Logger {
	logger := Logger{s, log.New(out, prefix, flag)}
	return &logger
}

type Logger struct {
	source *Source
	output *log.Logger
}

func (logger *Logger) Debug(content string) {
	if (*(*logger).source).logLevel == "debug" {
		(*(*logger).source).mutex.Lock()
		(*logger).output.Println(content)
		(*(*logger).source).mutex.Unlock()
	}
}

func (logger *Logger) Info(content string) {
	if (*(*logger).source).logLevel != "warning" {
		(*(*logger).source).mutex.Lock()
		(*logger).output.Println(content)
		(*(*logger).source).mutex.Unlock()
	}
}

func (logger *Logger) Warning(content string) {
	(*(*logger).source).mutex.Lock()
	(*logger).output.Println(content)
	(*(*logger).source).mutex.Unlock()
}

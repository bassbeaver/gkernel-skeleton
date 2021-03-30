package logger_factory

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	"github.com/bassbeaver/logopher"
	"io"
	"os"
)

const (
	loggerFactoryServiceAlias = "LoggerFactory"
)

type LoggerFactory struct {
	logFilePath string
}

func (f *LoggerFactory) CreateLogger() *logopher.Logger {
	logfile, logfileError := os.OpenFile(f.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if logfileError != nil {
		panic(logfileError)
	}

	// Wrapping stdout into anonymous structure with (embedded) io.Writer in order for the resulting stdoutWriter
	// to stop implementing the io.Closer interface to prevent stdout closing in middleware.RequestLoggerCloser
	stdoutWriter := &struct{ io.Writer }{Writer: os.Stdout}
	stdoutHandler := logopher.CreateStreamHandler(stdoutWriter, &logopher.SimpleFormatter{}, nil, nil, 1)

	fileHandler := logopher.CreateStreamHandler(logfile, &logopher.JsonFormatter{}, nil, nil, 4)

	logger := &logopher.Logger{}
	logger.SetHandlers([]logopher.HandlerInterface{stdoutHandler, fileHandler})

	return logger
}

//--------------------

func newLoggerFactory(path string) *LoggerFactory {
	return &LoggerFactory{logFilePath: path}
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(loggerFactoryServiceAlias, newLoggerFactory, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", loggerFactoryServiceAlias, err.Error()))
	}
}

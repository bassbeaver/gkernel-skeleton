package request_logger

import (
	"github.com/bassbeaver/gkernel/event_bus/event"
	"github.com/bassbeaver/gkernel/response"
	"github.com/bassbeaver/logopher"
	"io"
	"os"
	"path/filepath"
)

type Middleware struct {
	LogFilePath string
}

func (m *Middleware) CreateLogger(eventObj *event.RequestReceived) {
	logfile, logfileError := os.OpenFile(m.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
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

	eventObj.RequestContextAppend(RequestContextKey, logger)
}

func (m *Middleware) CloseLogger(eventObj *event.RequestTermination) {
	if _, isWsUpgrade := eventObj.GetResponse().(*response.WebsocketUpgradeResponse); isWsUpgrade {
		return
	}

	if loggerObj := eventObj.GetRequest().Context().Value(RequestContextKey); nil != loggerObj {
		logger := loggerObj.(*logopher.Logger)
		logger.ExportBufferedMessages()
		logger.CloseStreams()
	}
}

//--------------------

func NewMiddleware(LogFilePath string) *Middleware {
	absolutePath, absolutePathError := filepath.Abs(LogFilePath)
	if nil != absolutePathError {
		panic(absolutePathError)
	}

	return &Middleware{
		LogFilePath: absolutePath,
	}
}

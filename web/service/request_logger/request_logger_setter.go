package request_logger

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	"github.com/bassbeaver/gkernel/web/event_bus/event"
	"github.com/bassbeaver/gkernel/web/response"
	"github.com/bassbeaver/logopher"
	"gkernel-skeleton/common/service/logger_factory"
	"net/http"
)

const (
	MiddlewareServiceAlias = "RequestLoggerSetter"
	RequestContextKey      = "logger"
)

type RequestLoggerSetter struct {
	loggerFactory *logger_factory.LoggerFactory
}

func (m *RequestLoggerSetter) SetLoggerToRequestContext(eventObj *event.RequestReceived) {
	logger := m.loggerFactory.CreateLogger()
	eventObj.RequestContextAppend(RequestContextKey, logger)
}

func (m *RequestLoggerSetter) CloseLogger(eventObj *event.RequestTermination) {
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

func newRequestLoggerSetter(factory *logger_factory.LoggerFactory) *RequestLoggerSetter {
	return &RequestLoggerSetter{loggerFactory: factory}
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, newRequestLoggerSetter, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

func GetFromRequestContext(request *http.Request) *logopher.Logger {
	return request.Context().Value(RequestContextKey).(*logopher.Logger)
}

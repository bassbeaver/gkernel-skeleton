package request_logger

import (
	"fmt"
	"github.com/bassbeaver/gkernel"
	"github.com/bassbeaver/logopher"
	"net/http"
)

const (
	MiddlewareServiceAlias = "RequestLoggerSetter"
	RequestContextKey      = "logger"
)

func Register(kernelObj *gkernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, NewMiddleware, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

func GetFromRequestContext(request *http.Request) *logopher.Logger {
	return request.Context().Value(RequestContextKey).(*logopher.Logger)
}

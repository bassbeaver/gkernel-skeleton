package timeout_handler

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	kernelResponse "github.com/bassbeaver/gkernel/web/response"
	"net/http"
	"strconv"
)

const (
	serviceAlias = "TimeoutHandler"
)

type TimeoutHandler struct {
}

func (th *TimeoutHandler) HandleTimeout() kernelResponse.Response {
	response := kernelResponse.NewViewResponse("error/bad-request.gohtml")
	response.SetHttpStatus(http.StatusServiceUnavailable)
	response.SetData(map[string]string{
		"ErrorCode": strconv.Itoa(response.GetHttpStatus()),
		"Text":      "Request timeout",
	})

	return response
}

func (th *TimeoutHandler) PerformLoginHandleTimeout() kernelResponse.Response {
	response := kernelResponse.NewViewResponse("error/bad-request.gohtml")
	response.SetHttpStatus(http.StatusServiceUnavailable)
	response.SetData(map[string]string{
		"ErrorCode": strconv.Itoa(response.GetHttpStatus()),
		"Text":      "Login process timeout",
	})

	return response
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(serviceAlias, func() *TimeoutHandler { return &TimeoutHandler{} }, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", serviceAlias, err.Error()))
	}
}

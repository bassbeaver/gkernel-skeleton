package request_size_validator

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	webEvent "github.com/bassbeaver/gkernel/web/event_bus/event"
	"github.com/bassbeaver/gkernel/web/response"
	"net/http"
	"strconv"
	"strings"
)

const (
	MiddlewareServiceAlias = "RequestSizeValidator"
	maxRequestBodySize     = 8 << 20 // 8 Mb
	requestBodySizeError   = "http: request body too large"
)

type Validator struct{}

func (g *Validator) ValidateBodySize(eventObj *webEvent.RequestReceived) {
	checkError := checkRequestBody(eventObj.GetResponseWriter(), eventObj.GetRequest())

	if nil == checkError {
		return
	}

	errorResponse := response.NewViewResponse("error/bad_request.gohtml")

	switch checkError.(type) {
	case *requestEntityToLarge:
		errorResponse.SetHttpStatus(http.StatusRequestEntityTooLarge)
		errorResponse.SetData(map[string]string{
			"ErrorCode": strconv.Itoa(http.StatusRequestEntityTooLarge),
			"Text":      "Request body too large",
		})
	default:
		errorResponse.SetHttpStatus(http.StatusBadRequest)
		errorResponse.SetData(map[string]string{
			"ErrorCode": strconv.Itoa(http.StatusBadRequest),
			"Text":      checkError.Error(),
		})
	}

	eventObj.SetResponse(errorResponse)
}

func newValidator() *Validator {
	return &Validator{}
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, newValidator, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

func checkRequestBody(responseWriter http.ResponseWriter, request *http.Request) error {
	method := request.Method
	if !("POST" == method || "PUT" == method || "DELETE" == method || "PATCH" == method) {
		return nil
	}

	request.Body = http.MaxBytesReader(responseWriter, request.Body, maxRequestBodySize)

	var formError error
	if strings.Contains(request.Header.Get("Content-Type"), "multipart/form-data") {
		formError = request.ParseMultipartForm(maxRequestBodySize)
	} else {
		formError = request.ParseForm()
	}

	if nil != formError {
		if requestBodySizeError == formError.Error() {
			return &requestEntityToLarge{error: formError}
		}

		return formError
	}

	return nil
}

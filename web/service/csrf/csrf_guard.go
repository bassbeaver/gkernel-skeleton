package csrf

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	"github.com/bassbeaver/gkernel/web/event_bus/event"
	"github.com/bassbeaver/gkernel/web/response"
	sessionService "gkernel-skeleton/web/service/session"
	"net/http"
	"strconv"
)

const (
	middlewareServiceAlias = "CsrfGuard"
	tokenSessionKey        = "csrf_token"
	formInputName          = "csrf_token"
)

type Guard struct{}

func (g *Guard) Set(eventObj *event.RequestReceived) {
	sessionObj := sessionService.GetFromRequestContext(eventObj.GetRequest())

	csrfToken := sessionObj.GetString(tokenSessionKey)
	if "" != csrfToken {
		return
	}

	randomBytes := make([]byte, sha256.Size224)
	_, randomErr := rand.Read(randomBytes)
	if nil != randomErr {
		panic(randomErr)
	}

	csrfTokenBytes := sha256.Sum256(randomBytes)
	csrfToken = fmt.Sprintf("%x", csrfTokenBytes)

	sessionObj.Set(tokenSessionKey, csrfToken)
}

func (g *Guard) Check(eventObj *event.RequestReceived) {
	method := eventObj.GetRequest().Method
	if !("POST" == method || "PUT" == method || "DELETE" == method || "PATCH" == method) {
		return
	}

	request := eventObj.GetRequest()

	formError := request.ParseForm()
	if nil != formError {
		errorResponse := response.NewViewResponse("error/bad-request.gohtml")
		errorResponse.SetHttpStatus(http.StatusBadRequest)
		errorResponse.SetData(map[string]string{
			"ErrorCode": strconv.Itoa(http.StatusBadRequest),
			"Text":      formError.Error(),
		})
		eventObj.SetResponse(errorResponse)

		return
	}

	formToken := request.PostFormValue(formInputName)
	sessionToken := GetTokenFromRequestContext(request)
	if formToken != sessionToken {
		errorResponse := response.NewViewResponse("error/bad-request.gohtml")
		errorResponse.SetHttpStatus(http.StatusBadRequest)
		errorResponse.SetData(map[string]string{
			"ErrorCode": strconv.Itoa(http.StatusBadRequest),
			"Text":      "Invalid CSRF token",
		})
		eventObj.SetResponse(errorResponse)
	}
}

func newCsrfGuard() *Guard {
	return &Guard{}
}

func GetTokenFromRequestContext(request *http.Request) string {
	return sessionService.GetFromRequestContext(request).GetString(tokenSessionKey)
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(middlewareServiceAlias, newCsrfGuard, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", middlewareServiceAlias, err.Error()))
	}
}

package csrf

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/bassbeaver/gkernel"
	"github.com/bassbeaver/gkernel/event_bus/event"
	"github.com/bassbeaver/gkernel/response"
	sessionService "gkernel-skeleton/service/session"
	"net/http"
	"strconv"
)

const (
	MiddlewareServiceAlias = "CsrfGuard"
	TokenSessionKey        = "csrf_token"
	FormInputName          = "csrf_token"
)

type Guard struct{}

func (g *Guard) Set(eventObj *event.RequestReceived) {
	sessionObj := sessionService.GetFromRequestContext(eventObj.GetRequest())

	csrfToken := sessionObj.GetString(TokenSessionKey)
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

	sessionObj.Set(TokenSessionKey, csrfToken)
}

func (g *Guard) Check(eventObj *event.RequestReceived) {
	method := eventObj.GetRequest().Method
	if !("POST" == method || "PUT" == method || "DELETE" == method || "PATCH" == method) {
		return
	}

	request := eventObj.GetRequest()

	formError := request.ParseForm()
	if nil != formError {
		errorResponse := response.NewViewResponse("cp/error/bad_request.gohtml")
		errorResponse.SetHttpStatus(http.StatusBadRequest)
		errorResponse.SetData(map[string]string{
			"ErrorCode": strconv.Itoa(http.StatusBadRequest),
			"Text":      formError.Error(),
		})
		eventObj.SetResponse(errorResponse)

		return
	}

	formToken := request.PostFormValue(FormInputName)
	sessionToken := GetTokenFromRequestContext(request)
	if formToken != sessionToken {
		errorResponse := response.NewViewResponse("cp/error/bad_request.gohtml")
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
	return sessionService.GetFromRequestContext(request).GetString(TokenSessionKey)
}

func Register(kernelObj *gkernel.Kernel) {
	err := kernelObj.RegisterService(MiddlewareServiceAlias, newCsrfGuard, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", MiddlewareServiceAlias, err.Error()))
	}
}

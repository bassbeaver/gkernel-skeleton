package session

import (
	"fmt"
	webKernel "github.com/bassbeaver/gkernel/web"
	webEvent "github.com/bassbeaver/gkernel/web/event_bus/event"
	"github.com/kataras/go-sessions/v3"
	"net/http"
	"time"
)

const (
	sessionsMiddlewareService      = "SessionsMiddleware"
	sessionRequestContextKey       = "session"
	headersBufferRequestContextKey = "session_response_writer_buffer"
)

type Middleware struct {
	sessionsObj *sessions.Sessions
}

func (m *Middleware) RequestSessionStart(eventObj *webEvent.RequestReceived) {
	sessRW := newHeadersBuffer()

	session := m.sessionsObj.Start(sessRW, eventObj.GetRequest())

	eventObj.RequestContextAppend(sessionRequestContextKey, session)
	eventObj.RequestContextAppend(headersBufferRequestContextKey, sessRW)
}

func (m *Middleware) RequestSessionAddResponseHeader(eventObj *webEvent.ResponseBeforeSend) {
	sessionHeaders := eventObj.GetRequest().Context().Value(headersBufferRequestContextKey)
	if nil == sessionHeaders {
		return
	}

	// Copying headers from session headers to response headers
	for headerName, headerValues := range sessionHeaders.(*headersBuffer).Header() {
		for _, value := range headerValues {
			eventObj.GetResponse().GetHeaders().Add(headerName, value)
		}
	}
}

//--------------------

func newMiddleware(cookieId string, redisConnection sessions.Database) *Middleware {
	sessionsObj := sessions.New(sessions.Config{Cookie: cookieId, Expires: time.Hour * 24})
	sessionsObj.UseDatabase(redisConnection)

	return &Middleware{
		sessionsObj: sessionsObj,
	}
}

func Register(kernelObj *webKernel.Kernel) {
	err := kernelObj.RegisterService(sessionsMiddlewareService, newMiddleware, true)
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", sessionsMiddlewareService, err.Error()))
	}
}

func GetFromRequestContext(request *http.Request) *sessions.Session {
	return request.Context().Value(sessionRequestContextKey).(*sessions.Session)
}

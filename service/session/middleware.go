package session

import (
	"github.com/bassbeaver/gkernel/event_bus/event"
	"github.com/kataras/go-sessions"
	sessionsRedis "github.com/kataras/go-sessions/sessiondb/redis"
)

type Middleware struct {
	sessionsObj *sessions.Sessions
}

func (m *Middleware) InitRedisConnection(_ *event.ApplicationLaunched) {
	// Do just nothing, main goal is to build all necessary services
}

func (m *Middleware) CloseRedisConnection(eventObj *event.ApplicationTermination) {
	eventObj.GetContainer().GetByAlias(RedisConnectionServiceAlias).(*sessionsRedis.Database).Close()
}

func (m *Middleware) RequestSessionStart(eventObj *event.RequestReceived) {
	sessRW := newHeadersBuffer()

	session := m.sessionsObj.Start(sessRW, eventObj.GetRequest())

	eventObj.RequestContextAppend(SessionRequestContextKey, session)
	eventObj.RequestContextAppend(HeadersBufferRequestContextKey, sessRW)
}

func (m *Middleware) RequestSessionAddResponseHeader(eventObj *event.ResponseBeforeSend) {
	sessionHeaders := eventObj.GetRequest().Context().Value(HeadersBufferRequestContextKey)
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

func NewMiddleware(sessionsObj *sessions.Sessions) *Middleware {
	return &Middleware{
		sessionsObj: sessionsObj,
	}
}

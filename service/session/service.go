package session

import (
	"fmt"
	"github.com/bassbeaver/gkernel"
	"github.com/kataras/go-sessions"
	sessionsRedis "github.com/kataras/go-sessions/sessiondb/redis"
	sessionsRedisService "github.com/kataras/go-sessions/sessiondb/redis/service"
	"net/http"
	"time"
)

const (
	RedisConnectionServiceAlias = "SessionsRedisConnection"
	SessionsServiceAlias        = "Sessions"
	MiddlewareServiceAlias      = "SessionsMiddleware"

	SessionRequestContextKey       = "session"
	HeadersBufferRequestContextKey = "session_response_writer_buffer"
)

func newRedisConnection(
	redisAddress,
	redisPassword,
	redisDatabase string,
	redisMaxIdle,
	redisMaxActive,
	redisIdleTimeout int,
) *sessionsRedis.Database {
	redisConnection := sessionsRedis.New(sessionsRedisService.Config{
		Network:     sessionsRedisService.DefaultRedisNetwork,
		Addr:        redisAddress,
		Password:    redisPassword,
		Database:    redisDatabase,
		MaxIdle:     redisMaxIdle,
		MaxActive:   redisMaxActive,
		IdleTimeout: time.Duration(redisIdleTimeout) * time.Second,
		Prefix:      "",
	})

	return redisConnection
}

func newSessions(cookieId string, redisConnection sessions.Database) *sessions.Sessions {
	sessionsObj := sessions.New(sessions.Config{Cookie: cookieId, Expires: time.Hour * 24})
	sessionsObj.UseDatabase(redisConnection)

	return sessionsObj
}

func Register(kernelObj *gkernel.Kernel) {
	// Register services to Container
	var err error

	err = kernelObj.RegisterService(RedisConnectionServiceAlias, newRedisConnection, true)
	checkError(RedisConnectionServiceAlias, err)
	err = kernelObj.RegisterService(SessionsServiceAlias, newSessions, true)
	checkError(SessionsServiceAlias, err)
	err = kernelObj.RegisterService(MiddlewareServiceAlias, NewMiddleware, true)
	checkError(MiddlewareServiceAlias, err)
}

func GetFromRequestContext(request *http.Request) *sessions.Session {
	return request.Context().Value(SessionRequestContextKey).(*sessions.Session)
}

func checkError(serviceName string, err error) {
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", serviceName, err.Error()))
	}
}

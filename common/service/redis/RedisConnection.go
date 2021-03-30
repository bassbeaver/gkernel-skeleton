package redis

import (
	"fmt"
	"github.com/bassbeaver/gkernel/event_bus/event"
	sessionsRedis "github.com/kataras/go-sessions/v3/sessiondb/redis"
	sessionsRedisService "github.com/kataras/go-sessions/v3/sessiondb/redis/service"
	"gkernel-skeleton/common/lib"
	"time"
)

const (
	redisConnectionServiceAlias           = "RedisConnection"
	redisConnectionMiddlewareServiceAlias = "RedisConnectionMiddleware"
)

// Middleware used to init and close Redis connections
type ConnectionMiddleware struct{}

func (m *ConnectionMiddleware) InitRedisConnection(_ *event.ApplicationLaunched) {
	// Do just nothing, main goal is to build all necessary services
}

func (m *ConnectionMiddleware) CloseRedisConnection(eventObj *event.ApplicationTermination) {
	_ = eventObj.GetContainer().GetByAlias(redisConnectionServiceAlias).(*sessionsRedis.Database).Close()
}

//--------------------

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

func newMiddleware() *ConnectionMiddleware {
	return &ConnectionMiddleware{}
}

func Register(kernelObj lib.ServiceRegistrator) {
	err := kernelObj.RegisterService(redisConnectionServiceAlias, newRedisConnection, true)
	checkError(redisConnectionServiceAlias, err)

	err = kernelObj.RegisterService(redisConnectionMiddlewareServiceAlias, newMiddleware, true)
	checkError(redisConnectionMiddlewareServiceAlias, err)
}

func checkError(serviceName string, err error) {
	if nil != err {
		panic(fmt.Sprintf("failed to register %s service, error: %s", serviceName, err.Error()))
	}
}

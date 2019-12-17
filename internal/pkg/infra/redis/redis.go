package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

// Connection retorna conexão com Redis
type Connection interface {
	Get() redis.Conn
}

// Logger define interface para loggers
type Logger interface {
	Errorf(msg string, v ...interface{})
	Warnf(msg string, v ...interface{})
	Infof(msg string, v ...interface{})
	Debugf(msg string, v ...interface{})
	Tracef(msg string, v ...interface{})
}

// Redis representa pool de conexões Redis
type Redis struct {
	config config.Redis
	pool   *redis.Pool
	logger Logger
}

// NewRedis cria instância de Redis
func newRedis(config config.Redis, l Logger) (r *Redis, err error) {
	var rd = new(Redis)
	rd.config = config
	rd.logger = l

	var conTest redis.Conn
	conTest, err = redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return
	}
	defer conTest.Close()

	rd.pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
			if err != nil {
				l.Errorf("Conexão %v", err.Error())
			}
			return c, err
		},
	}

	r = rd
	r.logger.Infof("Conectado ao Redis\n")

	return
}

// Get retorna conexão com Redis do Pool
func (r *Redis) Get() (c redis.Conn) {
	c = r.pool.Get()
	return
}

// DBConnection representa conexão Redis em modo BD
type DBConnection Connection

// NewDBConnection retorna instância de DBConnection
func NewDBConnection(cfg config.RedisDB, l logging.LoggerRedisDB) (c DBConnection, err error) {
	c, err = newRedis(config.Redis(cfg), l)
	return
}

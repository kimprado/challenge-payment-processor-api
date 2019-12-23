package logging

import (
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	l "github.com/kimprado/sllog/pkg/logging"
)

// Logger para logar ROOT
type Logger struct {
	l.Logger
}

// LoggerProcessor para logar paymentprocessor.processor
type LoggerProcessor struct {
	l.Logger
}

// LoggerAPI para logar paymentprocessor.api
type LoggerAPI struct {
	l.Logger
}

// LoggerCardRepository para logar paymentprocessor.repository
type LoggerCardRepository struct {
	l.Logger
}

// LoggerRedisDB para logar infra.redis.db
type LoggerRedisDB struct {
	l.Logger
}

// LoggerWebServer para logar webserver
type LoggerWebServer struct {
	l.Logger
}

// NewLogger cria Logger ""(ROOT)
func NewLogger(configLevels config.LoggingLevels) (log Logger) {
	log = Logger{l.NewLogger("", configLevels)}
	return
}

// NewLoggerProcessor cria Logger "paymentprocessor.processor"
func NewLoggerProcessor(configLevels config.LoggingLevels) (log LoggerProcessor) {
	log = LoggerProcessor{l.NewLogger("paymentprocessor.processor", configLevels)}
	return
}

// NewLoggerAPI cria Logger "paymentprocessor.api"
func NewLoggerAPI(configLevels config.LoggingLevels) (log LoggerAPI) {
	log = LoggerAPI{l.NewLogger("paymentprocessor.api", configLevels)}
	return
}

// NewLoggerCardRepository cria Logger "paymentprocessor.repository"
func NewLoggerCardRepository(configLevels config.LoggingLevels) (log LoggerCardRepository) {
	log = LoggerCardRepository{l.NewLogger("paymentprocessor.repository", configLevels)}
	return
}

// NewRedisDB cria Logger "infra.redis.db"
func NewRedisDB(configLevels config.LoggingLevels) (log LoggerRedisDB) {
	log = LoggerRedisDB{l.NewLogger("infra.redis.db", configLevels)}
	return
}

// NewWebServer cria Logger "webserver"
func NewWebServer(configLevels config.LoggingLevels) (log LoggerWebServer) {
	log = LoggerWebServer{l.NewLogger("webserver", configLevels)}
	return
}

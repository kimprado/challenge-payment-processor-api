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

// LoggerHTTP para logar infra.http
type LoggerHTTP struct {
	l.Logger
}

// LoggerRedisDB para logar infra.redis.db
type LoggerRedisDB struct {
	l.Logger
}

// LoggerJWTFilter para logar infra.security.jwt
type LoggerJWTFilter struct {
	l.Logger
}

//LoggerWebConfigHTTPExporter para logar web.config.httpexporter
type LoggerWebConfigHTTPExporter struct {
	l.Logger
}

//LoggerWebInfoHTTPExporter para logar instrumentation.info.config
type LoggerWebInfoHTTPExporter struct {
	l.Logger
}

//LoggerWebVersionHTTPExporter para logar instrumentation.info.version
type LoggerWebVersionHTTPExporter struct {
	l.Logger
}

//LoggerMetricsRequestCounter para logar instrumentation.metrics.http.REQUEST_COUNTER
type LoggerMetricsRequestCounter struct {
	l.Logger
}

//LoggerMetricsRequestResponseTime para logar instrumentation.metrics.http.REQUEST_RESPONSE_TIME
type LoggerMetricsRequestResponseTime struct {
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

// NewLoggerHTTP cria Logger "infra.http"
func NewLoggerHTTP(configLevels config.LoggingLevels) (log LoggerHTTP) {
	log = LoggerHTTP{l.NewLogger("infra.http", configLevels)}
	return
}

// NewRedisDB cria Logger "infra.redis.db"
func NewRedisDB(configLevels config.LoggingLevels) (log LoggerRedisDB) {
	log = LoggerRedisDB{l.NewLogger("infra.redis.db", configLevels)}
	return
}

// NewLoggerJWTFilter cria Logger "infra.security.jwt"
func NewLoggerJWTFilter(configLevels config.LoggingLevels) (log LoggerJWTFilter) {
	log = LoggerJWTFilter{l.NewLogger("infra.security.jwt", configLevels)}
	return
}

//NewLoggerWebConfigHTTPExporter cria Logger "instrumentation.info.config"
func NewLoggerWebConfigHTTPExporter(configLevels config.LoggingLevels) (log LoggerWebConfigHTTPExporter) {
	log = LoggerWebConfigHTTPExporter{l.NewLogger("instrumentation.info.config", configLevels)}
	return
}

//NewLoggerWebInfoHTTPExporter cria Logger "instrumentation.info.info"
func NewLoggerWebInfoHTTPExporter(configLevels config.LoggingLevels) (log LoggerWebInfoHTTPExporter) {
	log = LoggerWebInfoHTTPExporter{l.NewLogger("instrumentation.info.info", configLevels)}
	return
}

//NewLoggerWebVersionHTTPExporter cria Logger "instrumentation.info.version"
func NewLoggerWebVersionHTTPExporter(configLevels config.LoggingLevels) (log LoggerWebVersionHTTPExporter) {
	log = LoggerWebVersionHTTPExporter{l.NewLogger("instrumentation.info.version", configLevels)}
	return
}

//NewMetricsRequestCounter cria Logger "instrumentation.metrics.http.REQUEST_COUNTER"
func NewMetricsRequestCounter(configLevels config.LoggingLevels) (log LoggerMetricsRequestCounter) {
	log = LoggerMetricsRequestCounter{l.NewLogger("instrumentation.metrics.http.REQUEST_COUNTER", configLevels)}
	return
}

//NewMetricsRequestResponseTime cria Logger "instrumentation.metrics.http.REQUEST_RESPONSE_TIME"
func NewMetricsRequestResponseTime(configLevels config.LoggingLevels) (log LoggerMetricsRequestResponseTime) {
	log = LoggerMetricsRequestResponseTime{l.NewLogger("instrumentation.metrics.http.REQUEST_RESPONSE_TIME", configLevels)}
	return
}

// NewWebServer cria Logger "webserver"
func NewWebServer(configLevels config.LoggingLevels) (log LoggerWebServer) {
	log = LoggerWebServer{l.NewLogger("webserver", configLevels)}
	return
}

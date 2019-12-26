package logging

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewFileAppender,

	NewLogger,
	NewLoggerProcessor,
	NewLoggerAPI,
	NewLoggerCardRepository,
	NewLoggerHTTP,
	NewRedisDB,
	NewLoggerWebConfigHTTPExporter,
	NewLoggerWebInfoHTTPExporter,
	NewLoggerWebVersionHTTPExporter,
	NewWebServer,
)

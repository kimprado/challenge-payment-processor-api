package logging

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewFileAppender,

	NewLogger,
	NewLoggerAPI,
	NewLoggerCardRepository,
	NewRedisDB,
	NewWebServer,
)

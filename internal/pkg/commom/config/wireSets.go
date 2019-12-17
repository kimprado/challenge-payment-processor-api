package config

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewConfig,
	NewLoggingLevels,
	NewRedisDB,
)

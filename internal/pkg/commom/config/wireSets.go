package config

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewConfig,
	NewLoggingLevels,
	NewRedisDB,
)

// AppSet define providers do pacote para inicialização da aplicação
var AppSet = wire.NewSet(
	NewLoggingLevels,
	NewRedisDB,
)

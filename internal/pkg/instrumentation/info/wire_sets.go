package info

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewApp,
	NewConfigExporterHTTP,
	NewAppInfoExporterHTTP,
	NewVersionExporterHTTP,
)

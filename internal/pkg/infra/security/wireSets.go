package security

import (
	"github.com/google/wire"
)

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewJWTFilter,
)

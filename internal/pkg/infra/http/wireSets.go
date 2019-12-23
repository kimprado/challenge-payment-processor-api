package http

import "github.com/google/wire"

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewHTTPService,
	// Define que a implementação Padão de RequestSender é Service
	wire.Bind(new(RequestSender), new(*Service)),
)

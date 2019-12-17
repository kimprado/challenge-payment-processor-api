package app

import (
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
	"github.com/challenge/payment-processor/internal/pkg/webserver"
	"github.com/google/wire"
)

// AppSet define providers do pacote
var AppSet = wire.NewSet(
	config.PkgSet,
	logging.PkgSet,
	redis.PkgSet,
	webserver.PkgSet,
	NewPaymentProcessorApp,
)

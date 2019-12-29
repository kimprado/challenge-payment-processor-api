package api

import (
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/http"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
	"github.com/challenge/payment-processor/internal/pkg/processor"
	"github.com/google/wire"
)

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewController,
)

var pkgSetConfigTest = wire.NewSet(
	newIntegrationConfigFile,
	config.PkgSet,
)

var pkgSetTest = wire.NewSet(
	PkgSet,
	config.NewLoggingLevels,
	config.NewRedisDB,
	logging.PkgSet,
	http.PkgSet,
	redis.PkgSet,
	processor.PkgSet,
)

type startWorkers struct{ ctrl *Controller }

func newStartWorkers(ctrl *Controller, sw *processor.StoneAcquirerWorkers, cw *processor.CieloAcquirerWorkers) (w *startWorkers) {
	return &startWorkers{ctrl}
}

func newIntegrationConfigFile() string {
	return "../../../../configs/config-integration.json"
}

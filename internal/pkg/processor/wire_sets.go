package processor

import (
	"github.com/google/wire"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
)

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewPaymentProcessorService,
	// Define que a implementação Padão de Processor é PaymentProcessorService
	wire.Bind(new(Processor), new(*PaymentProcessorService)),
	NewCardRepositoryRedis,
	// Define que a implementação Padão de CardRepositoryFinder é CardRepositoryRedis
	wire.Bind(new(CardRepositoryFinder), new(*CardRepositoryRedis)),
	NewAcquirerActors,
	// Define que a implementação Padão de AcquirerActorsSender é AcquirerActors
	wire.Bind(new(AcquirerActorsSender), new(*AcquirerActors)),
	// Define que a implementação Padão de AcquirerActorsResgister é AcquirerActors
	wire.Bind(new(AcquirerActorsResgister), new(*AcquirerActors)),
	NewActorsMap,
	NewStoneAcquirerWorkers,
	NewCieloAcquirerWorkers,
	NewAcquirerParameter,
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
	redis.PkgSet,
)

func newIntegrationConfigFile() string {
	return "../../../configs/config-integration.json"
}

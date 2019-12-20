package processor

import (
	"github.com/google/wire"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
)

// PkgSet define providers do pacote
var PkgSet = wire.NewSet(
	NewCardRepositoryRedis,
	// Define que a implementação Padão de CardRepositoryFinder é CardRepositoryRedis
	wire.Bind(new(CardRepositoryFinder), new(*CardRepositoryRedis)),
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

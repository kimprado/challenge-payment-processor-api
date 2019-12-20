// +build wireinject

package processor

import (
	"github.com/google/wire"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
)

// initializeConfigTest inicializa Configuration para testes
func initializeConfigTest() (config config.Configuration, err error) {
	panic(wire.Build(pkgSetConfigTest))
}

// initializeRedisTest inicializa DBConnection para testes
func initializeRedisTest(config config.Configuration) (dbc redis.DBConnection, err error) {
	panic(wire.Build(pkgSetTest))
}

// initializeCardRepositoryRedisTest inicializa CardRepositoryRedis para testes
func initializeCardRepositoryRedisTest(config config.Configuration) (rep *CardRepositoryRedis, err error) {
	panic(wire.Build(pkgSetTest))
}

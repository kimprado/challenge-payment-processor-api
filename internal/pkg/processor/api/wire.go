// +build wireinject

package api

import (
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
	"github.com/google/wire"
)

// initializeConfigTest inicializa Configuration para testes
func initializeConfigTest() (config config.Configuration, err error) {
	panic(wire.Build(pkgSetConfigTest))
}

// initializeRedisTest inicializa DBConnection para testes
func initializeRedisTest(config config.Configuration) (c redis.DBConnection, err error) {
	panic(wire.Build(pkgSetTest))
}

// initializeControllerWithDependenciesTest inicializa Controller e dependêicas para testes de integração
func initializeControllerWithDependenciesTest(config config.Configuration) (w *startWorkers, err error) {
	panic(wire.Build(wire.NewSet(pkgSetTest, newStartWorkers)))
}

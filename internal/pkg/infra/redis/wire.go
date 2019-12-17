// +build wireinject

package redis

import (
	"github.com/google/wire"
	"github.com/challenge/payment-processor/internal/pkg/commom/config"
)

// initializeConfigTest inicializa Configuration para testes
func initializeConfigTest() (config config.Configuration, err error) {
	panic(wire.Build(pkgSetConfigTest))
}

// initializeRedisTest inicializa DBConnection para testes
func initializeRedisTest(config config.Configuration) (c DBConnection, err error) {
	panic(wire.Build(pkgSetTest))
}

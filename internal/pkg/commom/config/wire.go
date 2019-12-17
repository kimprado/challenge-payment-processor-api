// +build wireinject

package config

import (
	"github.com/google/wire"
)

// initializeConfig inicializa Configuration para testes
func initializeConfigTest(path string) (config Configuration, err error) {
	panic(wire.Build(PkgSet))
}

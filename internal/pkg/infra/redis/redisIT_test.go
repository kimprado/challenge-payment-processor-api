// +build test integration

package redis

import (
	"testing"
)

func TestConectarDB(t *testing.T) {
	t.Parallel()

	var err error
	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	redis, err := initializeRedisTest(c)
	if err != nil {
		t.Errorf("Conexão banco de dados %v\n", err)
		return
	}

	cn := redis.Get()
	if cn == nil {
		t.Errorf("Connection não pode ser nula\n")
	}

}

// +build testenvvars

package config

import (
	"testing"
)

func TestNewConfigEnvVars(t *testing.T) {

	expect := struct {
		environment string
		serverPort  string
		redisDbHost string
		redisDbPort int
		logging     map[string]string
	}{
		environment: "test_ENV-VARS",
		serverPort:  "4033",
		redisDbHost: "host-env-test",
		redisDbPort: 6523,
		logging: map[string]string{
			"ROOT": "WARN-teste",
		},
	}

	var c Configuration
	var err error

	c, err = initializeConfigTest("../../../../configs/config-blank.json")

	if err != nil {
		t.Errorf("Erro ao carregar configurações %v", err)
		return
	}

	if expect.environment != c.Environment.Name {
		t.Errorf("Environment esperado %q é diferente de %q\n", expect.environment, c.Environment.Name)
	}
	if expect.serverPort != c.Server.Port {
		t.Errorf("serverPort esperado %q é diferente de %q\n", expect.serverPort, c.Server.Port)
	}
	if expect.redisDbHost != c.RedisDB.Host {
		t.Errorf("redisDbHost esperado %q é diferente de %q\n", expect.redisDbHost, c.RedisDB.Host)
	}
	if expect.redisDbPort != c.RedisDB.Port {
		t.Errorf("redisDbPort esperado %v é diferente de %v\n", expect.redisDbPort, c.RedisDB.Port)
	}

	for k, v := range expect.logging {
		z, ok := c.Logging.Level[k]
		if !ok {
			t.Errorf("Log de nível %q não encontrado na lista\n", k)
		}
		if ok && v != z {
			t.Errorf("Log esperado %q[%s] é diferente de %q[%s]\n", k, v, k, z)
		}
	}
}

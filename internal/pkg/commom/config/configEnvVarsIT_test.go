// +build test integration

package config

import (
	"os"
	"testing"
)

func TestNewConfigEnvVars(t *testing.T) {

	setUp := func() {
		os.Setenv("PROCESSOR_ENVIRONMENT_NAME", "test_ENV-VARS")
		os.Setenv("PROCESSOR_SERVER_PORT", "4033")
		os.Setenv("PROCESSOR_REDISDB_HOST", "host-env-test")
		os.Setenv("PROCESSOR_REDISDB_PORT", "6523")
		os.Setenv("PROCESSOR_STONEACQUIRER_URL", "http://local-test:8092/stone")
		os.Setenv("PROCESSOR_STONEACQUIRER_CONCURRENTWORKERS", "1000")
		os.Setenv("PROCESSOR_CIELOACQUIRER_URL", "http://local-test:8092/cielo")
		os.Setenv("PROCESSOR_CIELOACQUIRER_CONCURRENTWORKERS", "800")
		os.Setenv("PROCESSOR_SECURITY_JWTKEY", "app-jwt-key")
		os.Setenv("PROCESSOR_LOGGING_LEVEL", "ROOT: WARN-teste")
	}
	tearDown := func() {
		os.Setenv("PROCESSOR_ENVIRONMENT_NAME", "")
		os.Setenv("PROCESSOR_SERVER_PORT", "")
		os.Setenv("PROCESSOR_REDISDB_HOST", "")
		os.Setenv("PROCESSOR_REDISDB_PORT", "")
		os.Setenv("PROCESSOR_STONEACQUIRER_URL", "")
		os.Setenv("PROCESSOR_STONEACQUIRER_CONCURRENTWORKERS", "")
		os.Setenv("PROCESSOR_CIELOACQUIRER_URL", "")
		os.Setenv("PROCESSOR_CIELOACQUIRER_CONCURRENTWORKERS", "")
		os.Setenv("PROCESSOR_SECURITY_JWTKEY", "")
		os.Setenv("PROCESSOR_LOGGING_LEVEL", "")
	}
	setUp()
	defer tearDown()

	expect := struct {
		environment                    string
		serverPort                     string
		redisDbHost                    string
		redisDbPort                    int
		stoneAcquirerURL               string
		stoneAcquirerConcurrentWorkers int
		cieloAcquirerURL               string
		cieloAcquirerConcurrentWorkers int
		jwtKey                         string
		logging                        map[string]string
	}{
		environment:                    "test_ENV-VARS",
		serverPort:                     "4033",
		redisDbHost:                    "host-env-test",
		redisDbPort:                    6523,
		stoneAcquirerURL:               "http://local-test:8092/stone",
		stoneAcquirerConcurrentWorkers: 1000,
		cieloAcquirerURL:               "http://local-test:8092/cielo",
		cieloAcquirerConcurrentWorkers: 800,
		jwtKey:                         "app-jwt-key",
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
	if expect.stoneAcquirerURL != c.StoneAcquirer.URL {
		t.Errorf("stoneAcquirerURL esperado %v é diferente de %v\n", expect.stoneAcquirerURL, c.StoneAcquirer.URL)
	}
	if expect.stoneAcquirerConcurrentWorkers != c.StoneAcquirer.ConcurrentWorkers {
		t.Errorf("stoneAcquirerConcurrentWorkers esperado %v é diferente de %v\n", expect.stoneAcquirerConcurrentWorkers, c.StoneAcquirer.ConcurrentWorkers)
	}
	if expect.cieloAcquirerURL != c.CieloAcquirer.URL {
		t.Errorf("cieloAcquirerURL esperado %v é diferente de %v\n", expect.cieloAcquirerURL, c.CieloAcquirer.URL)
	}
	if expect.cieloAcquirerConcurrentWorkers != c.CieloAcquirer.ConcurrentWorkers {
		t.Errorf("cieloAcquirerConcurrentWorkers esperado %v é diferente de %v\n", expect.cieloAcquirerConcurrentWorkers, c.CieloAcquirer.ConcurrentWorkers)
	}
	if expect.jwtKey != c.Security.JWTKey {
		t.Errorf("jwtKey esperado %v é diferente de %v\n", expect.jwtKey, c.Security.JWTKey)
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

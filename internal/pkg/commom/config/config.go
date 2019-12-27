package config

import (
	"flag"
	"strconv"

	"github.com/jinzhu/configor"
)

// config possui configuracões da aplicação
var config = Configuration{}
var loaded bool

// Configuration - type
type Configuration struct {
	Environment struct {
		Name string `default:"dev"`
	}

	Server struct {
		Port string `default:"8081"`
	}

	RedisDB RedisDB

	StoneAcquirer struct {
		// URL string `default:"http://localhost:8092/stone"`
		URL               string `required:"true"`
		ConcurrentWorkers int    `required:"true"`
	}

	CieloAcquirer struct {
		// URL string `default:"http://localhost:8092/cielo"`
		URL               string `required:"true"`
		ConcurrentWorkers int    `required:"true"`
	}

	Security struct {
		Enable string `default:"true"`
		JWTKey string `required:"true"`
	}

	Logging struct {
		File  string `required:"false"`
		Level LoggingLevels
	}

	Metrics struct {
		Enable    bool   `default:"true"`
		Namespace string `default:"challenge"`
		Subsystem string `default:"payment_processor_api"`
	}
}

// Redis representa configuração de conexão Redis
type Redis struct {
	Host     string `default:"localhost"`
	Port     int    `default:"6379"`
	User     string `required:"false"`
	Password string `required:"false"`
	// Prefixo de todas chaves
	Prefix string `default:"processor"`
}

// RedisDB representa configuração Redis em modo DB
type RedisDB Redis

// NewRedisDB cria novo RedisDB
func NewRedisDB(c Configuration) (r RedisDB) {
	r = c.RedisDB
	return
}

// LoggingLevels representa loggers e seus respectivos níveis
type LoggingLevels map[string]string

// NewLoggingLevels cria novo LoggingLevels
func NewLoggingLevels(c Configuration) (ll LoggingLevels) {
	ll = c.Logging.Level
	return
}

// NewConfig -
func NewConfig(configLocationFile string) (c Configuration, err error) {
	var configLocation string
	if configLocationFile != "" {
		configLocation = configLocationFile
	} else {
		configLocation = loadFlags()
	}
	config, err = loadConfig(configLocation)
	if err != nil {
		return
	}
	if config.Security.Enable != "" {
		_, err = strconv.ParseBool(config.Security.Enable)
		if err != nil {
			return
		}
	}

	c = config
	return
}

func loadFlags() (configPath string) {
	cp := flag.String("config-location", "", "Caminho para arquivo de configuração")

	flag.Parse()

	configPath = *cp
	return
}

func loadConfig(configLocation string) (config Configuration, err error) {
	configApp := new(Configuration)

	cfg := configor.New(&configor.Config{
		ENVPrefix: "PROCESSOR",
	})

	if configLocation != "" {
		err = cfg.Load(configApp, configLocation)
	} else {
		err = cfg.Load(configApp)
	}

	if err != nil {
		return
	}
	config = *configApp
	return
}

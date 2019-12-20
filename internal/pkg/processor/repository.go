package processor

import (
	"encoding/json"
	"fmt"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
)

// CardRepositoryFinder consulta Cartao com token informado.
type CardRepositoryFinder interface {
	Find(token string) (c *Card, err error)
}

// CardRepositoryRedis implementa comportamentos de CardRepository
// com acesso a DB Redis.
type CardRepositoryRedis struct {
	redisClient redis.DBConnection
	redisCfg    config.RedisDB
	cfg         config.Configuration
	logger      logging.LoggerCardRepository
}

// NewCardRepositoryRedis cria instância de CardRepositoryRedis.
func NewCardRepositoryRedis(r redis.DBConnection, cr config.RedisDB, cfg config.Configuration, l logging.LoggerCardRepository) (c *CardRepositoryRedis) {
	c = new(CardRepositoryRedis)
	c.redisClient = r
	c.redisCfg = cr
	c.cfg = cfg
	c.logger = l
	return
}

// Find consulta cartão
func (cm *CardRepositoryRedis) Find(token string) (c *Card, err error) {

	con := cm.redisClient.Get()
	defer con.Close()

	var reply interface{}

	reply, err = con.Do("GET", fmt.Sprintf("%s:card:%s", cm.redisCfg.Prefix, token))
	if err != nil {
		return
	}
	if reply == nil {
		return
	}

	c = &Card{}
	err = json.Unmarshal(reply.([]byte), c)

	return
}

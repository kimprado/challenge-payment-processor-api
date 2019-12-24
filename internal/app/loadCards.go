package app

import (
	"encoding/json"
	"fmt"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/infra/redis"
	"github.com/challenge/payment-processor/internal/pkg/processor"
)

// DefaultCardsLoader carrega dados de cartões da aplicação
type DefaultCardsLoader struct {
	redisClient redis.DBConnection
	redisCfg    config.RedisDB
}

// NewDefaultCardsLoader cria instância de DefaultCardsLoader
func NewDefaultCardsLoader(r redis.DBConnection, cr config.RedisDB) (l *DefaultCardsLoader) {
	l = new(DefaultCardsLoader)
	l.redisClient = r
	l.redisCfg = cr
	return
}

func (l *DefaultCardsLoader) load() (err error) {

	loadCases := []struct {
		token string
		card  processor.Card
	}{
		{"xpto121a", processor.Card{Number: "121", CVV: "a"}},
		{"xpto122b", processor.Card{Number: "122", CVV: "b"}},
		{"xpto123c", processor.Card{Number: "123", CVV: "c"}},
	}

	con := l.redisClient.Get()
	defer con.Close()

	err = con.Send("MULTI")
	if err != nil {
		return
	}

	for _, lc := range loadCases {
		var cardJSON []byte
		cardJSON, err = json.Marshal(lc.card)
		con.Send("SET", fmt.Sprintf("%v:card:%v", l.redisCfg.Prefix, lc.token), string(cardJSON))
		if err != nil {
			return
		}
	}

	_, err = con.Do("EXEC")
	if err != nil {
		return
	}

	return
}

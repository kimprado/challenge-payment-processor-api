// +build test integration

package processor

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCards(t *testing.T) {
	t.Parallel()

	var err error

	err = setUpCards(t)
	if err != nil {
		t.Errorf("Erro ao preparar teste: %+v\n", err)
		return
	}

	c, err := initializeConfigTest()
	if err != nil {
		t.Errorf("Erro ao criar Configuração: %+v\n", err)
		return
	}

	c.RedisDB.Prefix = t.Name()

	repository, err := initializeCardRepositoryRedisTest(c)
	if err != nil {
		t.Errorf("Erro ao criar repository: %v\n", err)
		return
	}
	assert.NotNil(t, repository)

	testCases := []struct {
		token string
		card  *Card
	}{
		{"xpto121a", &Card{"121", "a"}},
		{"xpto122b", &Card{"122", "b"}},
		{"xpto123c", &Card{"123", "c"}},
		{"z000a", nil},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {

			result, err := repository.Find(tc.token)

			if err != nil {
				t.Errorf("Erro inesperado: %v", err)
				return
			}

			if tc.card == nil {
				assert.Nil(t, result)
				return
			}

			if invalid := assert.NotNil(t, result); invalid {
				return
			}
			assert.Equal(t, result.Number, tc.card.Number)
			assert.Equal(t, result.CVV, tc.card.CVV)

		})
	}

}

// setUpCards cria carga de dados para teste.
// Popula Redis com valores em nova chave para o teste.
// Nome da chave se baseia no nome do teste, o que permite
// executar testes de integração em paralelo :).
func setUpCards(t *testing.T) (err error) {

	loadCases := []struct {
		token string
		card  Card
	}{
		{"xpto121a", Card{"121", "a"}},
		{"xpto122b", Card{"122", "b"}},
		{"xpto123c", Card{"123", "c"}},
	}

	c, err := initializeConfigTest()
	if err != nil {
		return
	}

	redis, err := initializeRedisTest(c)
	if err != nil {
		return
	}

	con := redis.Get()
	defer con.Close()

	err = con.Send("MULTI")
	if err != nil {
		return
	}

	for _, lc := range loadCases {
		var cardJSON []byte
		cardJSON, err = json.Marshal(lc.card)
		con.Send("SET", fmt.Sprintf("%v:card:%v", t.Name(), lc.token), string(cardJSON))
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

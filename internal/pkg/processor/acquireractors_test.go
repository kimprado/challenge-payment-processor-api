// +build test unit

package processor

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRequestToActor(t *testing.T) {
	t.Parallel()

	var m ActorsMap
	var actors *AcquirerActors

	m = NewActorsMap()
	m["Stone"] = make(chan *AuthorizationRequest, 1)
	m["Cielo"] = make(chan *AuthorizationRequest, 1)
	actors = NewAcquirerActors(m)

	testCases := []struct {
		id          AcquirerID
		request     *AuthorizationRequest
		errExpected error
	}{
		{
			"Stone",
			nil,
			nil,
		},
		{
			"Stone",
			&AuthorizationRequest{},
			nil,
		},
		{
			"",
			nil,
			&AcquirerActorSendNotFoundError{},
		},
		{
			"",
			&AuthorizationRequest{},
			&AcquirerActorSendNotFoundError{},
		},
		{
			"Cielo",
			nil,
			nil,
		},
		{
			"Cielo",
			&AuthorizationRequest{},
			nil,
		},
		{
			"Rede",
			nil,
			&AcquirerActorSendNotFoundError{},
		},
		{
			"Rede",
			&AuthorizationRequest{},
			&AcquirerActorSendNotFoundError{},
		},
	}
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {

			err := actors.Send(tc.id, tc.request)

			if err != nil && tc.errExpected == nil {
				t.Errorf("Erro inesperado: %v", err)
				return
			}

			got := errors.Is(err, tc.errExpected)

			if got && tc.errExpected != nil {
				return
			}

			if !got && tc.errExpected != nil {
				t.Errorf("Esperado erro %v, mas obtido %v", tc.errExpected, err)
				return
			}

			var resultRequest *AuthorizationRequest

			resultRequest = <-m[tc.id]

			assert.Equal(t, tc.request, resultRequest)
			if tc.request == nil {
				assert.Nil(t, resultRequest)
			}

		})
	}

}

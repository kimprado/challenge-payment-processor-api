// +build test unit

package processor

import (
	"errors"
	"testing"
)

func TestSendRequestToActor(t *testing.T) {
	t.Parallel()

	actors := NewAcquirerActors()

	testCases := []struct {
		id          AcquirerID
		request     *AuthorizationRequest
		errExpected error
	}{
		{
			"stone",
			nil,
			nil,
		},
		{
			"teste-inexistente",
			nil,
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

		})
	}

}

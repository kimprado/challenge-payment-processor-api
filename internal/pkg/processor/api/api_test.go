// +build test unit

package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestProcessTransaction(t *testing.T) {
	var err error

	l := logging.NewLoggerAPI(nil)
	handler := NewController(nil, l)
	router := httprouter.New()
	router.POST("/transactions/", handler.Process)

	req, err := http.NewRequest("POST", "/transactions/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

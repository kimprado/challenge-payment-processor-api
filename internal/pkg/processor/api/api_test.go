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
	l := logging.NewLoggerAPI(nil)
	handler := NewController(nil, l)
	router := httprouter.New()
	router.POST("/transactions/", handler.Process)

	req, _ := http.NewRequest("POST", "/transactions/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

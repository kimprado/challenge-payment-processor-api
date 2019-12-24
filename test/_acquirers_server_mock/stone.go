package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/challenge/payment-processor/internal/pkg/processor"
)

func stoneHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delay)

	var err error

	var dto processor.AcquirerTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dto.Total > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Valor inválido"}`)
	}

	fmt.Fprintf(w, `{"message":"Transação autorizada"}`)
}

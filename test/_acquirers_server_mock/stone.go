package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func stoneHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delay)

	var err error

	logRequest(w, r)

	var dto AcquirerTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dto.Total > 1000 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Valor inválido"}`)
		return
	}

	if dto.Holder == "João Antônio" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Transação não permitida para o portador"}`)
		return
	}

	if dto.Installments > 12 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Não aprovado"}`)
		return
	}

	fmt.Fprintf(w, `{"message":"Transação autorizada"}`)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func cieloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delay)

	var err error

	logRequest(w, r)

	var dto AcquirerTransactionDTO
	err = json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dto.Total > 500 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Valor inválido"}`)
		return
	}

	if dto.Holder == "João Antônio" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Transação não permitida para o portador"}`)
		return
	}

	if dto.Installments > 6 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"Não aprovado"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Transação autorizada"}`)
}

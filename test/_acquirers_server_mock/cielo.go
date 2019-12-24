package main

import (
	"fmt"
	"net/http"
	"time"
)

func cieloHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delay)
	fmt.Fprintf(w, `{"message":"Transação autorizada"}`)
}

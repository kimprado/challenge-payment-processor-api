package main

import (
	"fmt"
	"net/http"
	"time"
)

func stoneHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(delay)
	fmt.Fprintf(w, "stone acquirer ")
}
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var port = os.Getenv("ACQUIRERS_PORT")
var strDelay = os.Getenv("ACQUIRERS_DELAY")
var delay time.Duration

func init() {

	log.Printf("Iniciando servidor")

	if port == "" {
		port = "8092"
	}

	if strDelay == "" {
		strDelay = "100"
	}

	intDelay, err := strconv.Atoi(strDelay)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	delay = (time.Duration)(intDelay) * time.Millisecond
}

func main() {

	http.HandleFunc("/stone", stoneHandler)
	http.HandleFunc("/cielo", cieloHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

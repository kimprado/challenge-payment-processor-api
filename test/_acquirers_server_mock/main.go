package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var port = os.Getenv("ACQUIRERS_PORT")
var strDelay = os.Getenv("ACQUIRERS_DELAY")
var delay time.Duration
var logging = os.Getenv("ACQUIRERS_LOGGING")

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
	log.Printf("Servidor rodando na porta :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func logRequest(w http.ResponseWriter, r *http.Request) {
	if logging == "DEBUG" {
		buf, bodyErr := ioutil.ReadAll(r.Body)
		if bodyErr != nil {
			http.Error(w, "Falha ao recuperar Request.Body", http.StatusInternalServerError)
		} else {
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			log.Printf("\n%s\n", rdr1)
			rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr2
		}
	}
}

package web

import (
	"encoding/json"
	"net/http"
)

// HTTPResponse representa resposta ao cliente HTTP
type HTTPResponse struct {
	writer  http.ResponseWriter
	code    int
	value   interface{}
	message error
}

// NewHTTPResponse cria instância de HTTPResponse
func NewHTTPResponse(r http.ResponseWriter, c int, value interface{}, e error) (hr *HTTPResponse) {
	hr = new(HTTPResponse)
	hr.writer = r
	hr.code = c
	hr.value = value
	hr.message = e
	return
}

// WriteJSON envia resposta HTTP.
// Converte value em JSON, define status code, e
// enventualmente mensagem de erro.
func (hr *HTTPResponse) WriteJSON() (err error) {

	hr.writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	hr.writer.WriteHeader(int(hr.code))

	if hr.message != nil {
		json.NewEncoder(hr.writer).Encode(hr.message)
		return

	}
	json.NewEncoder(hr.writer).Encode(hr.value)

	return
}

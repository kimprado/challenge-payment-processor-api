package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var errHTTPStatusBadRequest = errors.New("HTTP Bad Request")
var errHTTPServerError = errors.New("HTTP Server Error")

// RequestSender representa serviço de envio de requisições HTTP
type RequestSender interface {
	Send(url string, body interface{}, response interface{}) (err error)
}

// Service implementa utilitário para fazer requisições HTTP.
type Service struct {
	client *http.Client
}

// NewHTTPService cria instância de HTTPService
func NewHTTPService() (h *Service) {
	h = new(Service)
	h.client = &http.Client{
		Timeout: time.Second * 30,
	}
	return
}

// Send envia requisições POST, para URL com parâmetros informados
func (h *Service) Send(url string, body interface{}, response interface{}) (err error) {

	b, err := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return errHTTPStatusBadRequest
	}

	if resp.StatusCode >= 500 {
		return errHTTPServerError
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)

	return
}

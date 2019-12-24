package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

// StatusBadRequestError representa erro 400
type StatusBadRequestError struct {
	Message string
	Err     error
}

func newStatusBadRequestError(err error) (e *StatusBadRequestError) {
	e = new(StatusBadRequestError)
	e.Message = "HTTP Bad Request"
	e.Err = err
	return
}

func (e *StatusBadRequestError) Error() string {
	return e.Message
}

// Is informa se target == e. Verifica se e é do tipo
// StatusBadRequestError.
func (e *StatusBadRequestError) Is(target error) bool {
	switch target.(type) {
	case *StatusBadRequestError:
		return true
	default:
		return false
	}
}

// Unwrap retorna erro interno
func (e *StatusBadRequestError) Unwrap() error { return e.Err }

// ServerError representa erro 5xx
type ServerError struct {
	Message string
	Err     error
}

func newServerError(err error) (e *ServerError) {
	e = new(ServerError)
	e.Message = "HTTP Bad Request"
	e.Err = err
	return
}

func (e *ServerError) Error() string {
	return e.Message
}

// Is informa se target == e. Verifica se e é do tipo
// ServerError.
func (e *ServerError) Is(target error) bool {
	switch target.(type) {
	case *ServerError:
		return true
	default:
		return false
	}
}

// Unwrap retorna erro interno
func (e *ServerError) Unwrap() error { return e.Err }

var errHTTPServerError = errors.New("HTTP Server Error")

// Error representa erro em requisição HTTP
type Error struct {
	URL     string
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return e.Message
}

// RequestSender representa serviço de envio de requisições HTTP
type RequestSender interface {
	Send(url string, body interface{}, response interface{}) (err error)
}

// Service implementa utilitário para fazer requisições HTTP.
type Service struct {
	client *http.Client
	logger logging.LoggerHTTP
}

// NewHTTPService cria instância de HTTPService
func NewHTTPService(l logging.LoggerHTTP) (h *Service) {
	h = new(Service)
	h.client = &http.Client{
		Timeout: time.Second * 30,
	}
	h.logger = l
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

	errMessage := struct {
		Message string `json:"message"`
	}{}

	if resp.StatusCode == http.StatusBadRequest {
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&errMessage)
		err = newStatusBadRequestError(&Error{URL: url, Message: errMessage.Message, Code: resp.StatusCode})
		return
	}

	if resp.StatusCode >= 500 {
		err = newServerError(&Error{Code: resp.StatusCode})
		return
	}

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		err = &Error{Message: "Request error", Err: &Error{URL: url, Message: errMessage.Message, Code: resp.StatusCode}}
		return
	}

	decoder := json.NewDecoder(resp.Body)
	errDecode := decoder.Decode(&response)
	if errDecode != nil {
		h.logger.Errorf("Decode: %v\n", errDecode)
	}

	return
}

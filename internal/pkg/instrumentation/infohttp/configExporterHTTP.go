package infohttp

import (
	"encoding/json"
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

//ConfigExporterHTTP expõe informações da configuração da aplicação
type ConfigExporterHTTP struct {
	config config.Configuration
	logger logging.LoggerWebConfigHTTPExporter
}

//NewConfigExporterHTTP cria instância de ConfigExporterHTTP
func NewConfigExporterHTTP(c config.Configuration, l logging.LoggerWebConfigHTTPExporter) (e *ConfigExporterHTTP) {
	e = &ConfigExporterHTTP{
		config: c,
		logger: l,
	}
	return
}

// Serve envia configurações
func (e *ConfigExporterHTTP) Serve(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := json.NewEncoder(res).Encode(e.config)
	if err != nil {
		e.logger.Errorf("Processar config - %s\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao carregar config")
		return
	}
}

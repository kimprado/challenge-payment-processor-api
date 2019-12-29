package info

import (
	"fmt"
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

//VersionExporterHTTP expõe versão da aplicação
type VersionExporterHTTP struct {
	version string
	logger  logging.LoggerWebVersionHTTPExporter
}

//NewVersionExporterHTTP cria instância de VersionExporterHTTP
func NewVersionExporterHTTP(a App, l logging.LoggerWebVersionHTTPExporter) (e *VersionExporterHTTP) {
	e = &VersionExporterHTTP{
		version: a.Version.Version,
		logger:  l,
	}
	return
}

// Serve envia versão
func (e *VersionExporterHTTP) Serve(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	fmt.Fprintf(res, e.version)
}

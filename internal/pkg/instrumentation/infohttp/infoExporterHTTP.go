package infohttp

import (
	"encoding/json"
	"net/http"

	"runtime/debug"

	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
)

// InfoExporterHTTP expõe informações aplicação
type InfoExporterHTTP struct {
	info   App
	logger logging.LoggerWebInfoHTTPExporter
}

// NewInfoExporterHTTP cria instância de InfoExporterHTTP
func NewInfoExporterHTTP(a App, l logging.LoggerWebInfoHTTPExporter) (e *InfoExporterHTTP) {
	e = &InfoExporterHTTP{
		info:   a,
		logger: l,
	}
	return
}

// Serve envia informações da aplicação
func (e *InfoExporterHTTP) Serve(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var info *debug.BuildInfo
	var ok bool
	if info, ok = debug.ReadBuildInfo(); ok {
		e.info.Version.VersaoModulo = info.Main.Version
	}

	err := json.NewEncoder(res).Encode(e.info)
	if err != nil {
		e.logger.Errorf("Processar info - %s\n", err)
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode("Falha ao carregar infos")
		return
	}
}

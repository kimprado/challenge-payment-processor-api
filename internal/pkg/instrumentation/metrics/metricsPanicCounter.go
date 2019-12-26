package metrics

import (
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//PanicCounter -
type PanicCounter struct {
	counter prometheus.Counter
	logger  logging.LoggerWebServer
}

//NewPanicCounter recebe dependências
func NewPanicCounter(c config.Configuration, l logging.LoggerWebServer) (m *PanicCounter) {
	m = new(PanicCounter)
	m.counter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: c.Metrics.Namespace,
		Subsystem: c.Metrics.Subsystem,
		Name:      "http_req_fail_panic_total",
		Help:      "Total de requests em que ocorreu Panic",
	})
	m.logger = l

	return
}

//Handle coleta métricas para Prometheus
func (m *PanicCounter) Handle(w http.ResponseWriter, r *http.Request, err interface{}) {
	m.logger.Errorf("Request error: %s - %s ", r.URL.Path, err)
	m.counter.Inc()
	w.WriteHeader(http.StatusInternalServerError)
}

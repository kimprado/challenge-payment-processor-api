package metrics

import (
	"fmt"
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//RequestCounter -
type RequestCounter struct {
	counter *prometheus.CounterVec
	logger  logging.LoggerMetricsRequestCounter
}

//NewRequestCounter recebe dependências
func NewRequestCounter(c config.Configuration, l logging.LoggerMetricsRequestCounter) (m *RequestCounter) {
	m = new(RequestCounter)
	m.counter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: c.Metrics.Namespace,
		Subsystem: c.Metrics.Subsystem,
		Name:      "http_req_processed_total",
		Help:      "Total de requests recebidos",
	},
		[]string{
			"code", "method", "url",
		},
	)
	m.logger = l

	return
}

//RequestCounter envia métricas para Prometheus
func (m *RequestCounter) metricsRequestCounter(code int, method string, url string) {
	if !m.logger.IsTraceEnabled() {
		m.logger.Debugf("HTTP.Status %v - %q - %q\n", code, method, url)
	}
	m.counter.WithLabelValues(fmt.Sprint(code), method, url).Inc()
}

type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
	return &metricsResponseWriter{w, http.StatusOK}
}

func (mrw *metricsResponseWriter) WriteHeader(code int) {
	mrw.statusCode = code
	mrw.ResponseWriter.WriteHeader(code)
}

//Wrap coleta métricas para Prometheus
func (m *RequestCounter) Wrap(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if m.logger.IsTraceEnabled() {
			m.logger.Tracef("--> %s %s", req.Method, req.URL.Path)
		}

		mrw := newMetricsResponseWriter(w)
		wrappedHandler.ServeHTTP(mrw, req)

		statusCode := mrw.statusCode
		m.metricsRequestCounter(statusCode, req.Method, req.URL.Path)
		if m.logger.IsTraceEnabled() {
			m.logger.Tracef("<-- %d %s", statusCode, http.StatusText(statusCode))
		}
	})
}

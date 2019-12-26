package metrics

import (
	"net/http"
	"time"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

//ReqResponseTime -
type ReqResponseTime struct {
	hist   prometheus.Histogram
	logger logging.LoggerMetricsRequestResponseTime
}

//NewReqResponseTime recebe deendências
func NewReqResponseTime(c config.Configuration, l logging.LoggerMetricsRequestResponseTime) (m *ReqResponseTime) {
	m = new(ReqResponseTime)

	m.hist = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: c.Metrics.Namespace,
		Subsystem: c.Metrics.Subsystem,
		Name:      "http_req_response_time",
		Help:      "Total de requests em que ocorreu Panic",
	})
	m.logger = l

	return
}

//Handle coleta métricas para Prometheus
func (m *ReqResponseTime) Handle(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		start := time.Now()
		next(res, req, params)
		duration := time.Since(start)
		sec := float64(duration.Nanoseconds()) / float64(1000000000)
		m.logger.Debugf("%vs (%v)\n", sec, duration)
		m.hist.Observe(float64(sec))
	}
}

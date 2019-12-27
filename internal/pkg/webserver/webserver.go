package webserver

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/rileyr/middleware"

	cfg "github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	"github.com/challenge/payment-processor/internal/pkg/infra/security"
	"github.com/challenge/payment-processor/internal/pkg/instrumentation/info"
	"github.com/challenge/payment-processor/internal/pkg/instrumentation/metrics"
	"github.com/challenge/payment-processor/internal/pkg/processor/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var portNumber = regexp.MustCompile("^\\d{1,5}$")

var router *httprouter.Router

// WebServer representa servidor web que atende requisições http
type WebServer struct {
	*ParamWebServer
	home *Home
}

// ParamWebServer encapsula parâmetros de WebServer
type ParamWebServer struct {
	ctrl            *api.Controller
	configExporter  *info.ConfigExporterHTTP
	infoExporter    *info.AppInfoExporterHTTP
	versionExporter *info.VersionExporterHTTP
	reqResponseTime *metrics.ReqResponseTime
	requestcounter  *metrics.RequestCounter
	panicCounter    *metrics.PanicCounter
	jwt             *security.JWTFilter
	config          cfg.Configuration
	logger          logging.LoggerWebServer
}

// NewParamWebServer cria referência ParamWebServer
func NewParamWebServer(c *api.Controller, exporterHTTP *info.ConfigExporterHTTP, infoExporterHTTP *info.AppInfoExporterHTTP, versionExporterHTTP *info.VersionExporterHTTP, mrt *metrics.ReqResponseTime, mrc *metrics.RequestCounter, mpc *metrics.PanicCounter, jwt *security.JWTFilter, config cfg.Configuration, l logging.LoggerWebServer) (p *ParamWebServer) {
	p = new(ParamWebServer)
	p.ctrl = c
	p.configExporter = exporterHTTP
	p.infoExporter = infoExporterHTTP
	p.versionExporter = versionExporterHTTP
	p.reqResponseTime = mrt
	p.requestcounter = mrc
	p.panicCounter = mpc
	p.jwt = jwt
	p.config = config
	p.logger = l
	return
}

// NewWebServer cria referência WebServer
func NewWebServer(p *ParamWebServer) (w *WebServer) {
	w = new(WebServer)
	w.ParamWebServer = p
	w.home = NewHome()

	return
}

func serveHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	res.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PUT, DELETE, OPTIONS, HEAD, PATCH")
	router.ServeHTTP(res, req)
}

// Start é responsável por inicializar o servidor http
func (ws *WebServer) Start() {

	router = httprouter.New()

	mid := middleware.NewStack()

	if ws.config.Metrics.Enable {
		mid.Use(ws.reqResponseTime.Handle)
	}

	if ws.config.Security.Enable {
		mid.Use(ws.jwt.Handle)
	}

	router.GET("/", ws.home.Serve)
	router.POST("/transactions", mid.Wrap(ws.ctrl.Process))

	var defaultHandler http.Handler
	defaultHandler = http.HandlerFunc(serveHTTP)
	if ws.config.Metrics.Enable {
		defaultHandler = ws.requestcounter.Wrap(defaultHandler)
	}

	http.Handle("/", defaultHandler)
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/config", http.HandlerFunc(ws.configExporter.Serve))
	http.Handle("/info", http.HandlerFunc(ws.infoExporter.Serve))
	http.Handle("/version", http.HandlerFunc(ws.versionExporter.Serve))

	if ws.config.Metrics.Enable {
		router.PanicHandler = ws.panicCounter.Handle
	}

	var serverPort = ws.config.Server.Port
	if portNumber.MatchString(serverPort) {
		serverPort = "0.0.0.0:" + serverPort
	}

	var srv http.Server
	srv.Addr = serverPort

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint
		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	ws.logger.Infof("Servidor rodando na porta %v\n", serverPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		ws.logger.Errorf("Erro ao subir o servidor na porta %v - %s\n", serverPort, err)
		return
	}
	<-idleConnsClosed
}

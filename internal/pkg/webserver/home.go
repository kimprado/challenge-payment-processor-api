package webserver

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Home renderiza página da API
type Home struct {
}

// NewHome Cria instância de Home
func NewHome() (i *Home) {
	return &Home{}
}

// Serve responde conteúdo da página inicial
func (e Home) Serve(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	html := `
	<html>
	<body>
		<h1>
			Payment Processor API
		</h1>
		<br/>
		Links Úteis:
		<ul>
			<li>
				<a href='/docs'>Documentação</a> - Documentação interativa da API em formato Swagger.
			</li>
			<li>
				<a href='./info'>info</a> - Informações sobre a aplicação. JSON com versão, revisão, nome, data da compilação, etc.
			</li>
			<li>
				<a href='./version'>version</a> - Versão da aplicação(TAG). Informação preenchida na compilação ao gerar binário com 'make build' ou pelo CI.
			</li>
			<li>
				<a href='./config'>config</a> - Configuração real usada em Runtime, não se restringe a, mas inicia com parametrização(Env/config-file).
			</li>
			<li>
				<a href='./metrics'>metrics</a> - Méricas/Informações no formato Prometheus, como: Versões das dependências, tempo de resposta, total de requests, Tempo de GC.
			</li>
			<li>
				<a href='http://localhost:3001/d/kKd-m3qiz/requisicoes-http-payment-processor-api?orgId=1&kiosk=tv'>Grafana</a> - Visualização no Grafana das métricas coletadas pelo Prometheus.
			</li>
			<li>
				<a href='http://localhost:9090/targets'>Prometheus</a> - Coleta métricas da aplicação e permite consultas sobre as mesmas.
			</li>
		</ul>
	</body>
	</html>
	`

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, html)

}

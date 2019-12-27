package security

import (
	"net/http"

	"github.com/challenge/payment-processor/internal/pkg/commom/config"
	"github.com/challenge/payment-processor/internal/pkg/commom/logging"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

// JWTFilter aplica validação de segurança.
// Verifica token JWT da requisição.
type JWTFilter struct {
	jwtKey []byte
	logger logging.LoggerJWTFilter
}

// NewJWTFilter cria instância de JWTFilter
func NewJWTFilter(c config.Configuration, l logging.LoggerJWTFilter) (j *JWTFilter) {
	j = new(JWTFilter)
	j.jwtKey = []byte(c.Security.JWTKey)
	j.logger = l
	return
}

// Handle verifica existência e validade do token JWT.
func (j *JWTFilter) Handle(next httprouter.Handle) httprouter.Handle {

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {

		authHeader := req.Header.Get("Authorization")

		if len(authHeader) < 8 {
			j.logger.Debugf("Token inválido\n")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := authHeader[7:]

		tk, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return j.jwtKey, nil
		})

		if err != nil {
			j.logger.Debugf("Erro ao processar token: %v\n", err)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, _ := tk.Claims.(jwt.MapClaims)

		if !tk.Valid {
			j.logger.Debugf("Acesso negado para usuário: %v\n", claims["sub"])
			res.WriteHeader(http.StatusForbidden)
			return
		}

		app := claims["aud"]

		if app != "payment-processor-api" {
			j.logger.Debugf("Acesso negado para usuário: %v\n", claims["sub"])
			res.WriteHeader(http.StatusForbidden)
			return
		}

		next(res, req, params)
	}
}

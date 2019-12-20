package http

// RequestSender representa serviço de envio de requisições HTTP
type RequestSender interface {
	Send()
}

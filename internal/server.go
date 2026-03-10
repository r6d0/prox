package internal

const tcpProtocol = "tcp"
const digitBase = 10
const localhost = "127.0.0.1"
const httpsPort = "443"

// The abstraction of the proxy server.
type Prox interface {
	// The function starts the proxy server.
	Start() error
	// The function stops the proxy server.
	Stop() error
}

// The function creates a new instance of the proxy server.
func NewProx(config *ProxConfig) (Prox, error) {
	switch config.Protocol {
	case HTTP:
		return NewProxHttp(config)
	case HTTPS:
		return NewProxHttp(config)
	default:
		return nil, ErrUnsupportedProtocol
	}
}

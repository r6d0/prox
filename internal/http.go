package internal

import (
	"context"
	"io"
	"maps"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const XForwardedForHeader = "X-Forwarded-For"

var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

// HTTP proxy server.
type HttpProx struct {
	client *http.Client
	server *http.Server
	config *ProxConfig
	logger *LoggerWrapper
}

func (prox *HttpProx) Start() error {
	if prox.config.Protocol == HTTP {
		return prox.server.ListenAndServe()
	} else {
		tls := prox.config.Tls
		return prox.server.ListenAndServeTLS(tls.CertFile, tls.KeyFile)
	}
}

func (prox *HttpProx) Stop() error {
	return prox.server.Shutdown(context.Background())
}

func (prox *HttpProx) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	start := time.Now()
	status, err := prox.handle(wrt, req)

	logger := prox.logger
	logger.Debug(req.RemoteAddr, req.Host, status, time.Since(start), err)
}

func (prox *HttpProx) handle(wrt http.ResponseWriter, req *http.Request) (int, error) {
	if req.Method == http.MethodConnect {
		return prox.handleHttpConnect(wrt, req)
	} else {
		defer req.Body.Close()

		req.RequestURI = ""
		removeHopByHopHeaders(req.Header)

		config := prox.config
		if config.Request.Forwarded {
			req.Header.Add(XForwardedForHeader, req.RemoteAddr)
		}

		status := http.StatusInternalServerError
		res, err := prox.client.Do(req)
		if res != nil {
			defer res.Body.Close()

			removeHopByHopHeaders(res.Header)

			maps.Copy(wrt.Header(), res.Header)
			io.Copy(wrt, res.Body)
		}
		wrt.WriteHeader(status)
		return status, err
	}
}

func (prox *HttpProx) handleHttpConnect(wrt http.ResponseWriter, req *http.Request) (int, error) {
	status := http.StatusBadGateway
	host := req.Host

	_, _, err := net.SplitHostPort(host)
	if err != nil {
		host = net.JoinHostPort(host, httpsPort)
	}

	var target net.Conn
	timeout := time.Duration(prox.config.Request.Timeout)
	if target, err = net.DialTimeout(tcpProtocol, host, timeout); err == nil {
		defer target.Close()

		if hjr, ok := wrt.(http.Hijacker); !ok {
			err = http.ErrHijacked
			wrt.WriteHeader(status)
		} else {
			status = http.StatusOK
			wrt.WriteHeader(status)

			var source net.Conn
			if source, _, err = hjr.Hijack(); source != nil {
				defer source.Close()
			}

			group := sync.WaitGroup{}
			group.Add(2) // 2 is goroutines count.

			go func() { defer group.Done(); io.Copy(target, source) }()
			go func() { defer group.Done(); io.Copy(source, target) }()
			group.Wait()
		}
	}
	return status, err
}

// The function creates new instance of HTTP proxy server.
func NewProxHttp(config *ProxConfig) (Prox, error) {
	prox := &HttpProx{config: config, logger: NewLoggerWrapper(config)}

	port := strconv.FormatUint(uint64(config.Port), digitBase)
	addr := net.JoinHostPort(localhost, port)
	prox.server = &http.Server{Addr: addr, Handler: prox}
	prox.client = &http.Client{Timeout: time.Duration(config.Request.Timeout)}

	prox.logger.Info("Prox listens at [", addr, "].")
	return prox, nil
}

func removeHopByHopHeaders(headers http.Header) {
	for _, key := range hopByHopHeaders {
		headers.Del(key)
	}
}

package internal

import (
	"context"
	"io"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"os"
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
	logger *slog.Logger
}

func (prox *HttpProx) Start() error {
	return prox.server.ListenAndServe()
}

func (prox *HttpProx) Stop() error {
	return prox.server.Shutdown(context.Background())
}

func (prox *HttpProx) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	uri := req.RequestURI
	start := time.Now()
	status, err := prox.handle(wrt, req)

	prox.logger.Debug(
		"",
		"method", req.Method,
		"from", req.RemoteAddr,
		"to", uri,
		"status", status,
		"time", time.Since(start).Microseconds(),
		"err", err,
	)
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

			status = res.StatusCode
			wrt.WriteHeader(status)

			maps.Copy(wrt.Header(), res.Header)
			io.Copy(wrt, res.Body)
		} else {
			wrt.WriteHeader(status)
			prox.logger.Error(err.Error())
		}
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
	prox := &HttpProx{
		config: config,
		logger: slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.Log.Level})),
	}

	port := strconv.FormatUint(uint64(config.Port), digitBase)
	prox.server = &http.Server{Addr: portSeparator + port, Handler: prox}
	prox.client = &http.Client{Timeout: time.Duration(config.Request.Timeout)}

	prox.logger.Info("Prox listens at", "port", port)
	return prox, nil
}

func removeHopByHopHeaders(headers http.Header) {
	for _, key := range hopByHopHeaders {
		headers.Del(key)
	}
}

package internal

import (
	"context"
	"io"
	"log/slog"
	"maps"
	"net"
	"net/http"
	"os"
	"prox/internal/rule"
	"strconv"
	"strings"
	"sync"
	"time"
)

const portSeparator = ":"
const tcpProtocol = "tcp"
const digitBase = 10
const localhost = "127.0.0.1"
const httpsPort = "443"
const XForwardedForHeader = "X-Forwarded-For"
const IGNORE_URI = "/"

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

type poolItem struct {
	Data []byte
}

// HTTP proxy server.
type Prox struct {
	rules  []rule.Rule
	pool   *sync.Pool
	client *http.Client
	server *http.Server
	config *ProxConfig
	logger *slog.Logger
}

func (prox *Prox) Start() error {
	return prox.server.ListenAndServe()
}

func (prox *Prox) Stop() error {
	return prox.server.Shutdown(context.Background())
}

func (prox *Prox) ServeHTTP(wrt http.ResponseWriter, req *http.Request) {
	uri := req.RequestURI
	if !strings.HasPrefix(uri, IGNORE_URI) {
		prox.logger.Debug(
			"request_begin",
			"method", req.Method,
			"from", req.RemoteAddr,
			"to", uri,
		)

		// The rules checking.
		if len(prox.rules) > 0 {
			for _, rule := range prox.rules {
				if ok, code := rule.CheckHTTP(req); !ok {
					prox.logger.Debug(
						"request_not_allowed",
						"method", req.Method,
						"from", req.RemoteAddr,
						"to", uri,
						"code", code,
					)

					wrt.WriteHeader(code)
					return
				}
			}
			prox.logger.Debug(
				"request_allowed",
				"method", req.Method,
				"from", req.RemoteAddr,
				"to", uri,
			)
		}

		start := time.Now()
		status, err := prox.handle(wrt, req)

		prox.logger.Debug(
			"request_complete",
			"method", req.Method,
			"from", req.RemoteAddr,
			"to", uri,
			"status", status,
			"duration", time.Since(start).String(),
			"err", err,
		)
	}
}

func (prox *Prox) handle(wrt http.ResponseWriter, req *http.Request) (int, error) {
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

			// Copy all headers.
			maps.Copy(wrt.Header(), res.Header)

			// Copy request data.
			prox.copyBytes(req.Body, wrt)
		} else {
			wrt.WriteHeader(status)
			prox.logger.Error(err.Error())
		}
		return status, err
	}
}

func (prox *Prox) handleHttpConnect(wrt http.ResponseWriter, req *http.Request) (int, error) {
	defer req.Body.Close()

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

				if err == nil {
					group := sync.WaitGroup{}
					group.Add(2) // 2 is goroutines count.

					go func() {
						defer group.Done()

						var copyErr error
						for copyErr == nil || copyErr == io.EOF {
							_, copyErr = prox.copyBytes(target, source)
						}
					}()

					go func() {
						defer group.Done()

						var copyErr error
						for copyErr == nil || copyErr == io.EOF {
							_, copyErr = prox.copyBytes(source, target)
						}
					}()
					group.Wait()
				}
			}
		}
	}
	return status, err
}

func (prox *Prox) copyBytes(from io.Reader, to io.Writer) (int, error) {
	var read int
	var err error

	item := prox.pool.Get().(*poolItem)
	buffer := item.Data
	for true {
		if read, err = from.Read(buffer); read > 0 {
			to.Write(buffer[0:read])
		} else if read <= 0 || err != nil {
			break
		}
	}
	item.Data = buffer[:0]
	prox.pool.Put(item)

	return read, err
}

// The function creates new instance of HTTP proxy server.
func NewProx(config *ProxConfig) (*Prox, error) {
	rules, err := rule.NewRules(&config.Request.Rules)
	if err == nil {
		server := &Prox{
			rules:  rules,
			config: config,
			logger: createLogger(config),
			pool: &sync.Pool{
				New: func() any {
					return &poolItem{
						Data: make([]byte, config.Request.BufferSize),
					}
				},
			},
		}

		port := strconv.FormatUint(uint64(config.Port), digitBase)
		server.server = &http.Server{Addr: portSeparator + port, Handler: server}
		server.client = &http.Client{Timeout: time.Duration(config.Request.Timeout)}

		server.logger.Info("Prox listens at", "port", port)
		return server, nil
	}
	return nil, err
}

func removeHopByHopHeaders(headers http.Header) {
	for _, key := range hopByHopHeaders {
		headers.Del(key)
	}
}

func createLogger(config *ProxConfig) *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.Log.Level, AddSource: true}))
	return logger
}

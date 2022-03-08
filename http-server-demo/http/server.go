package http

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"gitee.com/lihp1603/utils/log"
)

type HttpApiServer struct {
	addr string
	mux  *HttpServerMux
	srv  *http.Server
}

func NewHttpApiServer(mux *HttpServerMux) *HttpApiServer {
	if nil == mux {
		mux = NewHttpServerMux()
	}
	return &HttpApiServer{mux: mux}
}

func (api *HttpApiServer) Start(addr string, tlsCertPath, tlsKeyPath string, enableTls bool) error {
	var err error = nil

	server := &http.Server{
		Addr:    addr,
		Handler: api.mux.mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second, // explicitly set no write timeout - see timeout handler.
	}
	api.addr = addr
	api.srv = server
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}
	if enableTls {
		log.Info("http server listen at %s enabled Tls", addr)
		// Start the HTTPS server
		if err = api.srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath); err != http.ErrServerClosed {
			log.Error("Error: failed to start server: %v", err)
		}
	} else {
		log.Info("http server listen at %s disabled Tls", addr)
		// Start the HTTP server
		if err = api.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Error: failed to start server: %v", err)
		}
	}
	return err
}

func (api *HttpApiServer) Stop() {
	if api.srv != nil {
		api.srv.Close()
	}
}

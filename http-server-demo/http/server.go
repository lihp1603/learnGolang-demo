package http

import (
	"gitee.com/lihp1603/utils/log"
	"net/http"
)

type HttpApiServer struct {
	addr         string
	srv          *http.Server
}



func NewHttpApiServer(addr string) *HttpApiServer {

}

type HttpServerMux struct {
	mux *http.ServeMux
}

func NewHttpServerMux() *HttpServerMux  {
	mux:=http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	return &HttpServerMux{mux:mux}
}

func (mux *HttpServerMux)RegisterHander(pattern string,method string, handler http.Handler)  {
	
}



func (api *HttpApiServer) Start(tlsCertPath, tlsKeyPath string, enableTls bool) error {
	var err error = nil
	go func() {
		if enableTls {
			log.Info("http server listen at %s enabled Tls", api.addr)
			// Start the HTTPS server
			if err = api.srv.ListenAndServeTLS(tlsCertPath, tlsKeyPath); err != http.ErrServerClosed {
				log.Error("Error: failed to start server: %v", err)
			}
		} else {
			log.Info("http server listen at %s disabled Tls", api.addr)
			// Start the HTTP server
			if err = api.srv.ListenAndServe(); err != http.ErrServerClosed {
				log.Error("Error: failed to start server: %v", err)
			}
		}
	}()

	return err
}

func (api *HttpApiServer) Close() {
	if api.srv != nil {
		api.srv.Close()
	}
}









package main

import (
	xhttp "development/http-server-demo/http"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitee.com/lihp1603/utils/log"
)

func onCreateNotifyConfig(w http.ResponseWriter, r *http.Request) {
	log.Info("request method:%s", r.Method)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello"))
	log.Info("write hello ok...")
}

func main() {
	mux := xhttp.NewHttpServerMux()
	mux.RegisterHander("/create", http.MethodGet, onCreateNotifyConfig)

	server := xhttp.NewHttpApiServer(mux)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		//wait
		<-c
		server.Stop()
		log.Info("stop server...")
	}()

	if err := server.Start(":9090", "", "", false); err != nil {
		log.Error("%s", err.Error())
	}
	return
}

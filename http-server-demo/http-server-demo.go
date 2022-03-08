package main

import (
	xhttp "development/http-server-demo/http"
	"net/http"

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
	if err := server.Start(":8090", "", "", false); err != nil {
		log.Error("%s", err.Error())
	}
	return
}

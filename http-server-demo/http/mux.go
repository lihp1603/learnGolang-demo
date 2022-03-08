package http

import (
	"errors"
	"net/http"
	"runtime"

	"gitee.com/lihp1603/utils"
	"gitee.com/lihp1603/utils/log"
)

type HttpServerMux struct {
	mux *http.ServeMux
}

func NewHttpServerMux() *HttpServerMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.NotFoundHandler())
	return &HttpServerMux{mux: mux}
}

func (mux *HttpServerMux) RegisterHander(pattern string, method string, handler http.HandlerFunc) error {
	if "" == pattern {
		log.Error("the pattern is empty")
		return errors.New("the pattern is empty")
	}
	if method == "" {
		mux.mux.Handle(pattern, handlePanic(handler))
	} else {
		//请求来了以后,执行顺序依次为handlePanic->requireMethod->handler
		mux.mux.Handle(pattern, handlePanic(requireMethod(method, handler)))
	}
	return nil
}

// requireMethod returns an http.HandlerFunc that checks whether
// the method of a client request matches the expected method before
// calling f.
//
// If the client request method does not match the given method
// it returns an error and http.StatusMethodNotAllowed to the client.
func requireMethod(method string, f http.HandlerFunc) http.HandlerFunc {
	var ErrMethodNotAllowed = NewError(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))

	return func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Set("Accept", method)
			sendHttpResError(w, ErrMethodNotAllowed)
			return
		}
		f(w, r)
	}
}

func handlePanic(f http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				switch err.(type) {
				case runtime.Error: // 运行时错误
					log.Error("runtime error: %s", err)
				default: // 非运行时错误
					log.Error("error: %s", err)
				}
				utils.PanicTrace()
			}
		}()
		f(writer, request)
	}
}

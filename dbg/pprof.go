package dbg

import (
	"net/http"
	"net/http/pprof"
)

func PProfHandlerFunc(w http.ResponseWriter, r *http.Request) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.ServeHTTP(w, r)
}

func PProfHandler() http.Handler {
	return http.HandlerFunc(PProfHandlerFunc)
}

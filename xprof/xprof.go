package xprof

import (
	"net/http"
	"net/http/pprof"
)

// StartProfiling is start make it pprof up
func StartProfiling() {
	go func() {
		r := http.NewServeMux()
		// Register pprof handlers
		r.HandleFunc("/debug/pprof/", pprof.Index)
		r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		r.HandleFunc("/debug/pprof/profile", pprof.Profile)
		r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		r.HandleFunc("/debug/pprof/trace", pprof.Trace)

		http.ListenAndServe(":6557", r)
	}()

}

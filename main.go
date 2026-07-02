package main

import (
	"fmt"
    "log"
    "net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
    const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

    mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("/healthz", handlerKickstart)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("/reset", apiCfg.handlerReset)

    srv := &http.Server{
        Addr :   ":" + port,
        Handler: mux,
	}

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    err := srv.ListenAndServe()
    if err != nil {
        log.Fatal(err)
    }
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	metric := cfg.fileserverHits.Load()
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", metric)))
}

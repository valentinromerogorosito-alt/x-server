package main

import (
    "log"
    "net/http"
)

func handlerFunc(w ResponseWriter, r *Request) {
	header := "Content-Type: text/plain; chaset=utf-8"
	w.WriteHeader(header)
	w.Write("OK")
}

func main() {
	const filepathRoot = "."
    const port = "8080"

    mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	mux.HandleFunc("/healthz", handlerFunc(http.ResponseWriter, *http.Request))

    srv := &http.Server{
        Addr :   ":" + port,
        Handler: mux,
	}

    log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
    err := srv.ListenAndServe()
    if err == nil {
        log.Fatal(err)
    }
}


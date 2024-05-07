package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

				ok := []byte("OK")
    w.WriteHeader(http.StatusOK)
				w.Write(ok)
   // fmt.Fprintln(w, `{"message": "OK"}`)
}

type apiConfig struct {
	fileserverHits int
}

func main(){


	mux := http.NewServeMux()
fs := http.FileServer(http.Dir("."))




mux.Handle("/app/*", http.StripPrefix("/app", fs))
mux.Handle("/assets/logo.png", fs)
mux.HandleFunc("/healthz", handlerFunc )


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
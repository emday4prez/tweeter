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

func serverHitsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")

				s:= fmt.Sprintf("Hits: %v", apiCfg.fileserverHits)
				ok := []byte(s)
    w.WriteHeader(http.StatusOK)
				w.Write(ok)
   // fmt.Fprintln(w, `{"message": "OK"}`)
}

func displayServerHitsHandler(w http.ResponseWriter, r *http.Request) {
	visits := apiCfg.fileserverHits

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    fmt.Fprintf(w, `
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
    `, visits)
			
    w.WriteHeader(http.StatusOK)
			

}

func resetHandler(w http.ResponseWriter, r *http.Request){
			apiCfg.resetServerHitCount()
    w.WriteHeader(http.StatusOK)
}

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	// Use a closure to wrap functionality
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++ // Increment the counter
		next.ServeHTTP(w, r)  // Call the next handler 
	})
}

func (c *apiConfig) resetServerHitCount( )  {
	c.fileserverHits = 0;
 
}

var apiCfg apiConfig





func main(){


	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	sfs := http.StripPrefix("/app", fs)


mux.Handle("/app/*", apiCfg.middlewareMetricsInc(sfs))
mux.Handle("/assets/logo.png", sfs)
mux.HandleFunc("GET /healthz", handlerFunc )
mux.HandleFunc("GET /metrics", serverHitsHandler )
mux.HandleFunc("GET /admin/metrics", displayServerHitsHandler )
mux.HandleFunc("/api/reset", resetHandler )


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
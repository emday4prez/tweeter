package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

type parameters struct {
    Body string `json:"body"`
}

type errorResponse struct {
    Error string `json:"error"`
}
type validResponse struct{
	Valid bool `json:"valid"`
}

type cleanedResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type chirpResponse struct {
	Id int `json:"id"`
	Body string `json:"body"`
}
 var uid = 1
var badWords = []string{"kerfuffle", "sharbert", "fornax"}
	

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request body"})
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Chirp is too long"})
		return
	}
// Use Fields to split by any whitespace
	words := strings.Fields(params.Body) 

	for i, word := range words {
		if containsIgnoreCase(badWords, word) {
			words[i] = "****"
		}
	}

	cleanedSentence := strings.Join(words, " ")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponse{Id: uid,Body: cleanedSentence })
	uid++
}

func containsIgnoreCase(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) { 
			return true
		}
	}
	return false
}

func chirpPostHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Invalid request body"})
		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "Chirp is too long"})
		return
	}
// Use Fields to split by any whitespace
	words := strings.Fields(params.Body) 

	for i, word := range words {
		if containsIgnoreCase(badWords, word) {
			words[i] = "****"
		}
	}

	cleanedSentence := strings.Join(words, " ")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponse{Id: uid,Body: cleanedSentence })
	uid++
 
}
 
func main(){
 
defer os.Remove("database.json")
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	sfs := http.StripPrefix("/app", fs)


mux.Handle("/app/*", apiCfg.middlewareMetricsInc(sfs))
mux.Handle("/assets/logo.png", sfs)
mux.HandleFunc("GET /healthz", handlerFunc )
mux.HandleFunc("GET /metrics", serverHitsHandler )
mux.HandleFunc("GET /admin/metrics", displayServerHitsHandler )
mux.HandleFunc("/api/reset", resetHandler )
mux.HandleFunc("/api/validate_chirp", validateHandler )
mux.HandleFunc("POST /api/chirps", chirpPostHandler )


	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}

}
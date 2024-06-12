package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/emday4prez/tweeter/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	DB             *database.DB
	jwtSecret      string
}


	func deleteDatabase(dbPath string) error {
		err := os.Remove(dbPath)
		if err != nil {
				return err
		}
		return nil
	}

func main() {
	godotenv.Load()
	const filepathRoot = "."
	const port = "8080"
dbg := flag.Bool("debug", false, "Enable debug mode")
flag.Parse()
	// Check if debug mode is enabled
	if *dbg {
		// Call the function to delete the database here.
		deleteDatabase("database.json")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}

jwtS := os.Getenv("JWT_SECRET")

	apiCfg := apiConfig{
		fileserverHits: 0,
		DB:             db,
		jwtSecret:						jwtS,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/*", fsHandler)

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /api/reset", apiCfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
 
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{id}", apiCfg.handlerChirpsGetById)

		mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
		mux.HandleFunc("PUT /api/users", apiCfg.handlerUsersUpdate)

		mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
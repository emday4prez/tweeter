package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
	

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
func (cfg *apiConfig) handlerChirpsGetById(w http.ResponseWriter, r *http.Request) {

 idString := r.PathValue("id")
	 id, err := strconv.Atoi(idString)
	dbChirp, err := cfg.DB.GetChirpById(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirps")
		return
	}

		c := Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		}

	respondWithJSON(w, http.StatusOK, c)
}


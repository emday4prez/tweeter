package main

//	"golang.org/x/crypto/bcrypt"
// 	"bcrypt"
import (
	"encoding/json"
	"net/http"
)

 

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {	
		Password string `json:password`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	
	validEmail, err := validateEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
 
	dbUsers, err := cfg.DB.GetUse()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}
 

	respondWithJSON(w, http.StatusCreated, User{
		ID:   user.ID,
		Email: user.Email,
	})
}



// func validateEmail(email string) (string, error) {
// 	const maxEmailLength = 140
// 	if len(email) > maxEmailLength {
// 		return "", errors.New("email is too long")
// 	}
// 	return email, nil
// }
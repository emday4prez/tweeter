package main

//	"golang.org/x/crypto/bcrypt"
//
import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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
 
	dbUser, err := cfg.DB.GetUserByEmail(validEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve User")
		return
	}
 
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(params.Password))
if err != nil {
   		respondWithError(w, http.StatusUnauthorized, "incorrect password")
} 

	respondWithJSON(w, http.StatusOK, User{
		ID:   dbUser.ID,
		Email: dbUser.Email,
	})
}



// func validateEmail(email string) (string, error) {
// 	const maxEmailLength = 140
// 	if len(email) > maxEmailLength {
// 		return "", errors.New("email is too long")
// 	}
// 	return email, nil
// }
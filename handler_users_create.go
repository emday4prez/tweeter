package main

//	"golang.org/x/crypto/bcrypt"
// 	"bcrypt"
import (
	"encoding/json"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID   int    `json:"id"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:password`
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
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password),14)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	user, err := cfg.DB.CreateUser(validEmail, string(hashedPassword))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:   user.ID,
		Email: user.Email,
	})
}

func validateEmail(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}
	return email, nil
}
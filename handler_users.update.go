package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ... (other imports)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        respondWithError(w, http.StatusUnauthorized, "Authorization header missing")
        return
    }

    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, jwt.ErrSignatureInvalid
        }
        return []byte(cfg.jwtSecret), nil
    })

    if err != nil || !token.Valid {
        respondWithError(w, http.StatusUnauthorized, "Invalid token")
        return
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        respondWithError(w, http.StatusInternalServerError, "Invalid token claims")
        return
    }

    userIDStr, ok := claims["sub"].(string)
    if !ok {
        respondWithError(w, http.StatusInternalServerError, "Invalid user ID in token")
        return
    }

    userID, err := strconv.Atoi(userIDStr)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Invalid user ID in token")
        return
    }

    type parameters struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
        return
    }

    // Hash the password before updating
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error hashing password")
        return
    }

    err = cfg.DB.UpdateUser(userID, params.Email, string(hashedPassword))
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't update user")
        return
    }

    updatedUser, err := cfg.DB.GetUserByEmail(params.Email)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve updated user")
        return
    }

    respondWithJSON(w, http.StatusOK, updatedUser)
}


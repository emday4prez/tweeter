package main

//	"golang.org/x/crypto/bcrypt"
//
import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

 

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {	
		Password string `json:"password"`
		Email string `json:"email"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
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
 fmt.Printf("\nLogin attempt: Email=%s \n", params.Email) 

	dbUser, err := cfg.DB.GetUserByEmail(validEmail)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve User")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password),[]byte(params.Password))
if err != nil {
	 fmt.Printf("Password comparison error: %v", err)
   		respondWithError(w, http.StatusUnauthorized, "incorrect password")
					return
} 

	expiresIn := 24 * time.Hour // default expiration
	if params.ExpiresInSeconds > 0 {
		potentialExpiresIn := time.Duration(params.ExpiresInSeconds) * time.Second
		if potentialExpiresIn <= 24*time.Hour {
			expiresIn = potentialExpiresIn
		}
	}
	expirationTime := time.Now().Add(expiresIn).Unix() 

    // Create JWT claims
    claims := jwt.MapClaims{
        "iss": "chirpy",    // Issuer
        "iat": time.Now(),    // Issued At
        "sub": strconv.Itoa(dbUser.ID),  // Subject (user ID)
        "exp": expirationTime, // Expiration (above)
    }

    // Create token
    token  := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				secretKey := cfg.jwtSecret

				// Sign token
    signedToken, err := token.SignedString([]byte(secretKey))
    if err != nil {
					fmt.Printf("error signing token:: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Error signing token")
        return
    }
// fmt.Printf("Login successful: UserID=%d, Token=%s", dbUser.ID, signedToken)
//     // Send JWT in the response header
//     w.Header().Set("Authorization", "Bearer " + signedToken)

// 	respondWithJSON(w, http.StatusOK, User{
// 		ID:   dbUser.ID,
// 		Email: dbUser.Email,
	 
// 	})

	// Prepare response
	response := map[string]interface{}{
		"id":    dbUser.ID,
		"email": dbUser.Email,
		"token": signedToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't encode response")
		return
	}

	fmt.Printf("User logged in: %s\n", dbUser.Email)
}
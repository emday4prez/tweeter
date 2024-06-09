package main

//	"golang.org/x/crypto/bcrypt"
//
import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

    // Create JWT claims
    claims := jwt.MapClaims{
        "iss": "your-auth-server",    // Issuer
        "sub": dbUser.ID,              // Subject (user ID)
        "exp": time.Now().Add(time.Hour * 24).Unix(), // Expiration (24 hours)
        // Add other relevant claims (e.g., roles) as needed
    }

    // Create token
    token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

    // Sign token
    signedToken, err := token.SignedString(cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error signing token")
        return
    }

    // Send JWT in the response header
    w.Header().Set("Authorization", "Bearer "+signedToken)

	respondWithJSON(w, http.StatusOK, User{
		ID:   dbUser.ID,
		Email: dbUser.Email,
	})
}
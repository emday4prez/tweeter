package main

//	"golang.org/x/crypto/bcrypt"
//
import (
	"crypto/rand"
	"encoding/hex"
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

type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	expiresIn := 1 * time.Hour // default expiration
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

    // Create JSON web token
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

 
	randomBytes := make([]byte, 32)

 retries := 3
    for i := 0; i < retries; i++ {
        _, err := rand.Read(randomBytes)
        if err != nil {
            // Log the error for analysis
            fmt.Printf("Error generating random bytes (attempt %d): %v", i+1, err)
            continue // Try again
        }
if retries < 1 {
respondWithError(w, 500, "Could not generate random bytes")
return
}
        // If no error, break out of the loop
        encodedStr := hex.EncodeToString(randomBytes)
        fmt.Println("Encoded random bytes:", encodedStr)	
								
								respondWithJSON(w, http.StatusOK, response{
									User: User{
									ID:    dbUser.ID,
									Email: dbUser.Email,
								},
									Token: signedToken,
									RefreshToken: encodedStr,
	})
        return // Exit the program successfully
    }

   

}
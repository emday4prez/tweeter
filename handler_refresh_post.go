package main

import (
	"fmt"
	"net/http"
	"strings"
)

type response struct {
		Token string `json:"token"`
	}

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request){


		    // check if the request method is not GET or HEAD
    if r.Method != http.MethodGet && r.Method != http.MethodHead {
        respondWithError(w, http.StatusMethodNotAllowed,"Method Not Allowed")
        return
    }

    // check if the Content-Length header is present and non-zero
    if r.ContentLength > 0 {
        respondWithError(w, http.StatusBadRequest,"Request body not allowed")
        return
    }
				// extract token
				reqToken := r.Header.Get("Authorization")
				splitToken := strings.Split(reqToken, "Bearer ")
				reqToken = splitToken[1]

				dbToken,err := cfg.DB.FindToken(reqToken)
				  if err != nil {
							respondWithError(w, 401, "cannot find token")
        return
    }
				res := response{
 Token: dbToken.Token,
}
fmt.Println(res)

respondWithJSON(w, 200, res)

}
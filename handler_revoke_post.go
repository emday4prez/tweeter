package main

import (
	"fmt"
	"net/http"
	"strings"
)

 

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request){

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

 
				res := response{
 
}
fmt.Println(res)

respondWithJSON(w, 200, res)

}
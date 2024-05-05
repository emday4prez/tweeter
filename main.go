package main

import (
	"fmt"
	"net/http"
)

func main(){
	mux := http.NewServeMux()
fs := http.FileServer(http.Dir("."))

mux.Handle("/", fs)

	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
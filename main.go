package main

import (
    "fmt"
    "log"
    "net/http"  

    "github.com/gorilla/mux"  
)

func main() {
    fmt.Println("started-service")

    r := mux.NewRouter()
    r.Handle("/upload", http.HandlerFunc(uploadHandler)).Methods("POST", "OPTIONS")
    r.Handle("/search", http.HandlerFunc(searchHandler)).Methods("GET", "OPTIONS")
    log.Fatal(http.ListenAndServe(":8080", r))
}


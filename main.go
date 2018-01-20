package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
        // vars := mux.Vars(r)
        // title := vars["title"]
        // page := vars["page"]

        fmt.Fprintf(w, "Hello World! :O")
    })

    http.ListenAndServe(":8000", r)
}
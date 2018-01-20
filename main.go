package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

// inspiration from:
// http://www.giantflyingsaucer.com/blog/?p=5635


func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World! :O")
}


func main() {

    r := mux.NewRouter()

    // serve static files
    s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)

    // main route
    r.Handle("/", http.FileServer(http.Dir("./static/")))

    // API routes
    r.HandleFunc("/hello", helloHandler)

    http.ListenAndServe(":8000", r)
}

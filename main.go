package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/mux"
)

// specifies either a singular int or list of ints that must be present together
type Requirement struct {
    needed []int
}

// says what options should be given if one of the requirements is met
type Path struct {
    options []int
    requirement []Requirement
}

type ActionText struct {
    action string
    text string
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    actionText map[int]ActionText
    paths map[int][]Path
}

// inspiration from:
// http://www.giantflyingsaucer.com/blog/?p=5635


func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World! :O")
}


func main() {

    r := mux.NewRouter()

    s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)
    // http.Handle("/", r)
    
    r.Handle("/", http.FileServer(http.Dir("./static/")))
    r.HandleFunc("/hello", helloHandler)

    http.ListenAndServe(":8000", r)
}
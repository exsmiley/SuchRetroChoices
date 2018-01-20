package main

import (
    // "fmt"
    // "net/http"

    // "github.com/gorilla/mux"
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
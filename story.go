package main

import (
)


type Element struct {
    action string
    text string
    points int
    together bool
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    elements map[int]Element
    paths map[int][]int
}
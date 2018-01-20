package main

import (
    "fmt"
    "log"
    "strings"
    "strconv"

    "gopkg.in/yaml.v2"
)


type Element struct {
    together bool
    action string
    text string
    points int
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    elements map[int]Element
    paths map[int][]int
}

var elementData = `
1: s; run; I want to run away; 500
2: t; talk; Wazzap dude; 100
`

var pathData = `
1: 2 3
2: 3 5 6 7
`


func LoadStory() Story {
    // create Story struct
    story := Story{}
    story.elements = make(map[int]Element)
    story.paths = make(map[int][]int)

    // load element data into a map
    elementMap := make(map[int]string)
    err := yaml.Unmarshal([]byte(elementData), &elementMap)
    if err != nil {
            log.Fatalf("error: %v", err)
    }

    for id, str := range elementMap {
        fields := strings.Split(str, ";")
        together, _ := strconv.ParseBool(strings.TrimSpace(fields[0]))
        points, _ := strconv.Atoi(strings.TrimSpace(fields[3]))
        element := Element{together, strings.TrimSpace(fields[1]), strings.TrimSpace(fields[2]), points}
        story.elements[id] = element
    }

    // load path data into a map
    pathMap := make(map[int]string)
    err = yaml.Unmarshal([]byte(pathData), &pathMap)
    if err != nil {
            log.Fatalf("error: %v", err)
    }

    for id, str := range pathMap {
        strPaths := strings.Split(str, " ")
        intPaths := []int{}

        for _, s := range strPaths {
            val, _ := strconv.Atoi(s)
            intPaths = append(intPaths, val)
        }

        story.paths[id] = intPaths
    }
    
    return story
}


func main() {
    story := LoadStory()
    fmt.Println(story)
}



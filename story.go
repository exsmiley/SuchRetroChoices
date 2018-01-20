package main

import (
    // "fmt"
    "log"
    "strings"
    "strconv"
    "io/ioutil"

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

func (story *Story) getElement(storyId int) Element {
    return story.elements[storyId]
}

func (story *Story) getPaths(storyId int) []int {
    return story.paths[storyId]
}

func (story *Story) hasEnded(storyId int) bool {
    return len(story.paths[storyId]) == 0
}


// loads the story object from the config files
func LoadStory() Story {
    // create Story struct
    story := Story{}
    story.elements = make(map[int]Element)
    story.paths = make(map[int][]int)

    // read files
    elementData, elementErr := ioutil.ReadFile("config/elements.yaml")
    if elementErr != nil {
        log.Printf("yamlFile.Get err   #%v ", elementErr)
    }
    pathData, pathErr := ioutil.ReadFile("config/paths.yaml")
    if pathErr != nil {
        log.Printf("yamlFile.Get err   #%v ", pathErr)
    }

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
        log.Println(id, intPaths)

        for _, s := range strPaths {
            val, _ := strconv.Atoi(s)
            intPaths = append(intPaths, val)
        }
        log.Println(id, intPaths)

        story.paths[id] = intPaths
    }
    
    return story
}

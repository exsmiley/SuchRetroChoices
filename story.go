package main

import (
    // "fmt"
    "log"
    "strings"
    "strconv"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type Path struct {
    normal []int
    help []int
    force []int
}


type Element struct {
    together bool
    action string
    text string
    points int
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    elements map[int]Element
    paths map[int]Path
}

func (story *Story) getElement(storyId int) Element {
    return story.elements[storyId]
}

func (story *Story) getPath(storyId int) Path {
    return story.paths[storyId]
}

func (story *Story) hasEnded(storyId int) bool {
    path := story.paths[storyId]
    return len(path.normal) == 0 && len(path.help) == 0 && len(path.force) == 0
}


// loads the story object from the config files
func LoadStory() Story {
    // create Story struct
    story := Story{}
    story.elements = make(map[int]Element)
    story.paths = make(map[int]Path)

    // read files
    elementData, elementErr := ioutil.ReadFile("config/elementstest.yaml")
    if elementErr != nil {
        log.Printf("yamlFile.Get err   #%v ", elementErr)
    }
    pathData, pathErr := ioutil.ReadFile("config/pathstest.yaml")
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
        typesOfPaths := [3][]int{}
        metaStrPaths := strings.Split(str, "|")
        for i, str2 := range metaStrPaths {
            strPaths := strings.Split(str2, " ")
            intPaths := []int{}

            for _, s := range strPaths {
                val, _ := strconv.Atoi(s)
                if val == 0 {
                    continue
                }
                intPaths = append(intPaths, val)
            }

            typesOfPaths[i] = intPaths
        }

        story.paths[id] = Path{typesOfPaths[0], typesOfPaths[1], typesOfPaths[2]}
    }
    
    return story
}

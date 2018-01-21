package main

import (
    // "fmt"
    "log"
    "strings"
    // "strconv"
    // "reflect"
    "io/ioutil"

    "gopkg.in/yaml.v2"
)

type Condition struct {
    requirements []string
    next string
}

type Choice struct {
    next string
    points int
    fallback string
    action string
}

type Special struct {
    event string
    result int
}

type Element struct {
    name string
    status bool
    text string
    special Special
    choices []Choice
    conditions []Condition
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    elements map[string]Element
}


func (story *Story) hasEnded(storyId string) bool {
    // path := story.paths[storyId]
    // return len(path.normal) == 0 && len(path.help) == 0 && len(path.force) == 0
    return true
}


// loads the story object from the config files
func LoadStory() Story {
    // create Story struct
    story := Story{}
    story.elements = make(map[string]Element)

    // read file
    elementData, elementErr := ioutil.ReadFile("config/elements.yaml")
    if elementErr != nil {
        log.Printf("yamlFile.Get err   #%v ", elementErr)
    }
    elementMap := make(map[string]map[string]interface{})
    err := yaml.Unmarshal([]byte(elementData), &elementMap)
    if err != nil {
            log.Fatalf("error: %v", err)
    }

    for k, val := range elementMap {
        el := Element{}
        el.name = k
        el.choices = []Choice{}
        el.conditions = []Condition{}
        for k2, val2 := range val {
            if k2 == "status" {
                // el.together = val2.(string) != "single"
            } else if k2 == "text" {
                // el.text = val2.(string)
            } else if k2 == "conditions" {
                
                val22 := val2.([]interface{})
                for _, v := range val22 {
                    condition := Condition{}
                    v2 := v.(map[interface{}]interface{})
                    for k4, v5 := range v2 {
                        k42 := k4.(string)
                        v52 := v5.(string)
                        if k42 == "condition" {
                            reqs := strings.Split(v52, "&")
                            for i := range reqs {
                                reqs[i] = strings.TrimSpace(reqs[i])
                            }
                            condition.requirements = reqs
                        } else if k42 == "next" {
                            condition.next = v52
                        } 
                        
                    }
                    el.conditions = append(el.conditions, condition)
                }
                
                
            } else if k2 == "choice" {
                val22 := val2.([]interface{})
                for _, v := range val22 {
                    choice := Choice{}
                    v2 := v.(map[interface{}]interface{})
                    for k4, v5 := range v2 {
                        k42 := k4.(string)
                        
                        if k42 == "points" {
                            v52 := v5.(int)
                            choice.points = v52
                        } else {
                            v52 := v5.(string)
                            if k42 == "next" {
                                choice.next = v52
                            } else if k42 == "fallback" {
                                choice.fallback = v52
                            } else if k42 == "action" {
                                choice.action = v52
                            }
                        }
                        
                    }
                    el.choices = append(el.choices, choice)
                }

            } else if k2 == "special" {
                special := Special{}
                v2 := val2.(map[interface{}]interface{})
                for k4, v5 := range v2 {

                    k42 := k4.(string)
                    if k42 == "event" {
                        v52 := v5.(string)
                        special.event = v52
                    } else {
                        v52 := v5.(int)
                        special.result = v52
                    }
                }
                el.special = special
            }
        }
        
        story.elements[k] = el
    }
    

    
    return story
}

func main() {
    log.Println("Loaded Story:", LoadStory())
}

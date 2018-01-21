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
    points int
}

type Choice struct {
    next string
    points int
    fallback string
    action string
    text string
    goback int
}

type Special struct {
    event string
    points int
    result string
}

type Element struct {
    name string
    status bool // true if single
    text string
    next string
    image string // url to an image
    special Special
    choices []Choice
    conditions []Condition
}

// object with text (id to text of the id) and paths (actions that can result from the specified id)
type Story struct {
    elements map[string]Element
}


func (story *Story) getText(storyId string) string {
    return story.elements[storyId].text
}

func (story *Story) getImage(storyId string) string {
    return story.elements[storyId].image
}

func (story *Story) getStatus(storyId string) string {
    return story.elements[storyId].text
}

func (story *Story) needsToWait(storyId string) bool {
    el := story.elements[storyId]
    log.Println("wait?", el)
    return !el.status || len(el.conditions) > 0
}

func (story *Story) triggersHelp(storyId string, action string) bool {
    el := story.elements[storyId]
    for _, choice := range el.choices {
        if choice.text == action && choice.action == "help" {
            return true
        }
    }
    return false
}

func (story *Story) triggersForce(storyId string, action string) bool {
    el := story.elements[storyId]
    for _, choice := range el.choices {
        if choice.text == action && choice.action == "force" {
            return true
        }
    }
    return false
}

// TODO future maybe filter by old actions
func (story *Story) getActions(storyId string, next string) []string {
    el := story.elements[storyId]
    actions := []string{}

    if next != "" {
        actions = append(actions, next)
    } else if len(el.choices) > 0 {
        for _, choice := range el.choices {
            actions = append(actions, choice.text)
        }
    } else if len(el.conditions) > 0 {
        action, _ := story.abortCondition(storyId)
        actions = append(actions, action)
    } else if el.next != "" {
        actions = append(actions, el.next)
    }
    return actions
}

func (story *Story) makeChoice(storyId string, action string) (string, int) {
    el := story.elements[storyId]
    for _, choice := range el.choices {
        if choice.text == action {
            return choice.next, choice.points
        }
    }

    for _, condition := range el.conditions {
        if condition.next == action {
            return condition.next, condition.points
        }
    }

    if el.next == action {
        return el.next, 0
    }

    return "", 0
}

func (story *Story) checkConditions(storyId string, action1 string, action2 string) string {
    el := story.elements[storyId]

    // this means that it was not actually a conditional spot, just needed some unity
    if len(el.conditions) == 0 {
        return storyId
    }

    for _, cond := range el.conditions {
        if (cond.requirements[0] == action1 && cond.requirements[1] == action2) || ( cond.requirements[1] == action1 && cond.requirements[0] == action2) {
            return cond.next
        }
    }

    abort, _ := story.abortCondition(storyId)

    return abort
}

func (story *Story) abortCondition(storyId string) (string, int) {
    el := story.elements[storyId]

    for _, cond := range el.conditions {
        if cond.requirements[0] != cond.requirements[1] {
            return cond.next, cond.points
        }
    }

    if el.next != "" {
        return el.next, 0
    }

    return "", 0
}

func (story *Story) hasEnded(storyId string) bool {
    el := story.elements[storyId]

    if el.special.event == "ENDING" {
        return true
    }

    return len(el.choices) == 0 && len(el.conditions) == 0 && el.special.event == "" && el.next == ""
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
        log.Println(val)
        for k2, val2 := range val {
            if k2 == "status" {
                el.status = val2.(string) == "single" || val2.(string) == ""
            } else if k2 == "text" {
                el.text = val2.(string)
            } else if k2 == "next" {
                el.next = val2.(string)
            } else if k2 == "image" {
                el.image = val2.(string)
            } else if k2 == "conditions" {
                
                val22 := val2.([]interface{})
                for _, v := range val22 {
                    condition := Condition{}
                    v2 := v.(map[interface{}]interface{})
                    for k4, v5 := range v2 {
                        log.Println(k4)
                        k42 := k4.(string)
                        if k42 == "special" {
                            continue
                        }

                        if k42 == "points" {
                            v52 := v5.(int)
                            condition.points = v52
                        } else {    
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
                        } else if k42 == "return" {
                            v52 := v5.(int)
                            choice.goback = v52
                        } else if k42 != "special" {
                            log.Println(k42)
                            v52 := v5.(string)
                            if k42 == "next" {
                                choice.next = v52
                            } else if k42 == "fallback" {
                                choice.fallback = v52
                            } else if k42 == "action" {
                                choice.action = v52
                            } else if k42 == "text" {
                                choice.text = v52
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
                    } else if k42 == "points" {
                        log.Println(v5)
                        v52 := v5.(int)
                        special.points = v52
                    } else if k42 == "result" {
                        v52 := v5.(string)
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

// func main() {
//     story := LoadStory()
//     log.Println(story)

// }

package main

import (
    // "log"
    "time"
    "crypto/sha256"
    "encoding/hex"
)


// taken from https://stackoverflow.com/questions/15323767/does-golang-have-if-x-in-construct-similar-to-python
func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// generates an id based on the current time
func genId() string {
    now := time.Now()
    h := sha256.New()
    h.Write([]byte(now.Format(time.RFC3339Nano)))
    return hex.EncodeToString(h.Sum(nil))
}
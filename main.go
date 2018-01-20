package main

import (
    "fmt"
    "net/http"
    "log"

    "github.com/googollee/go-socket.io"
    "github.com/gorilla/mux"
)

// inspiration from:
// http://www.giantflyingsaucer.com/blog/?p=5635


func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World! :O")
}


func main() {
    // start the mux router
    r := mux.NewRouter()

    // serve static files
    s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)

    // handle the socket
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }
    server.On("connection", socketHandler)
    server.On("error", socketErrorHandler)
    r.Handle("/socket.io/", server)

    // main route
    r.Handle("/", http.FileServer(http.Dir("./static/")))

    // API routes
    r.HandleFunc("/hello", helloHandler)

    log.Println("Started Server!")
    http.ListenAndServe(":8000", r)
}

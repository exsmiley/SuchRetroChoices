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
    // socket stuff
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }
    server.On("connection", func(so socketio.Socket) {
        log.Println("on connection")
        so.Join("chat")
        so.On("chat message", func(msg string) {
            log.Println("emit:", so.Emit("chat message", msg))
            so.BroadcastTo("chat", "chat message", msg)
        })
        so.On("disconnection", func() {
            log.Println("on disconnect")
        })
    })
    server.On("error", func(so socketio.Socket, err error) {
        log.Println("error:", err)
    })

    r := mux.NewRouter()

    // serve static files
    s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
    r.PathPrefix("/static/").Handler(s)

    // main route
    r.Handle("/", http.FileServer(http.Dir("./static/")))

    // API routes
    r.HandleFunc("/hello", helloHandler)
    r.Handle("/socket.io/", server)
    http.ListenAndServe(":8000", r)
}

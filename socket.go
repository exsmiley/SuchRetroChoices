package main

import (
    "log"

    "github.com/googollee/go-socket.io"
)


func socketHandler(so socketio.Socket) {
    log.Println("on connection")
    so.Join("chat")
    so.On("chat message", func(msg string) {
        log.Println(so.Emit("chat message", msg), " emit:", msg, )
        so.BroadcastTo("chat", "chat message", msg)
    })
    so.On("disconnection", func() {
        log.Println("on disconnect")
    })
}

func socketErrorHandler(so socketio.Socket, err error) {
    log.Println("error:", err)
}
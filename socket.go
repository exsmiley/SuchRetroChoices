package main

import (
    "log"

    "github.com/googollee/go-socket.io"
)


func metaSocketHandler(gm *GameMaster) func(so socketio.Socket) {
    return func(so socketio.Socket) {
        log.Println("on connection")

        // cookie, _ := so.Request().Cookie("yummy")
        so.Emit("chat message", "I don't know your name yet! :O")

        // first player needs to define player name
        so.On("name", func(playerName string) {
            if !gm.isInGame(playerName) {
                // TODO change to checking about the cookie
            }
        })

        so.Join("chat")
        so.On("chat message", func(msg string) {
            log.Println(so.Emit("chat message", msg), " emit:", msg, )
            so.BroadcastTo("chat", "chat message", msg)
        })
        so.On("disconnection", func() {
            log.Println("on disconnect")
        })
    }
}

func socketErrorHandler(so socketio.Socket, err error) {
    log.Println("error:", err)
}
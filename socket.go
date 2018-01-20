package main

import (
    "log"

    "github.com/googollee/go-socket.io"
)

func handleStartup() {
    // TODO handle logic for when the game starts and routing rooms and stuffs
}


func metaSocketHandler(gm *GameMaster) func(so socketio.Socket) {
    return func(so socketio.Socket) {
        log.Println("on connection")

        cookieObj, _ := so.Request().Cookie("yummy")
        cookie := cookieObj.Name

        // check if the user has set a name yet
        if gm.getName(cookie) == "" {
            so.Emit("ask name", "I don't know your name yet! :O")
        } else {
            so.Emit("chat message", "Welcome back " + gm.getName(cookie))
        }

        // first player needs to define player name
        so.On("name", func(playerName string) {
            // TODO ensure that len(playerName) > 0
            gm.setName(cookie, playerName)
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
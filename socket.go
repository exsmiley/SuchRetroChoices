package main

import (
    "log"
    "time"

    "github.com/googollee/go-socket.io"
)

/* Events
Client => Server
    name (str): name of the client
    join (str): cookie of the host to join
    action (str): action chosen by the client
    chat message (str): message to send to channel


Server => Client
    ask name (str): triggers the modal to ask for a name
    chat message (str): puts a message
    chat util (boolean): true if should display the chat
    room (object): data about who is hosting a game
    state (object): data about the current state of the game
*/

func handleStartup(gm *GameMaster, cookie string, so socketio.Socket) {
    // TODO handle logic for when the game starts and routing rooms and stuffs

    // check if the user has set a name yet
    if gm.getName(cookie) == "" {
        so.Emit("ask name", "I don't know your name yet! :O")
    } else {
        so.Emit("chat message", "Welcome back " + gm.getName(cookie))
    }

    // check if the user is in a room
    if gm.isInActiveGame(cookie) {
        so.Join(gm.getGame(cookie).Id)
        so.Emit("state", gm.getState(cookie))
    } else {
        so.Join("game room")
        gameId := gm.hostGame(cookie)
        so.Join(gameId)
        go reloadGameRoom(gm, cookie, so)
    }
}

// continually reloads
func reloadGameRoom(gm *GameMaster, cookie string, so socketio.Socket) {
    for !gm.isInActiveGame(cookie) {
        // TODO send room data
        so.Emit("room", gm.getRooms(cookie))
        time.Sleep(100 * time.Millisecond)
    }
}


func metaSocketHandler(gm *GameMaster) func(so socketio.Socket) {
    return func(so socketio.Socket) {
        cookieObj, _ := so.Request().Cookie("yummy")
        cookie := cookieObj.Value
        log.Println(cookie, "connected")

        handleStartup(gm, cookie, so)

        // first player needs to define player name
        so.On("name", func(playerName string) {
            // TODO ensure that len(playerName) > 0
            gm.setName(cookie, playerName)
        })

        so.On("join", func(host string) {
            // leave the game that player host alone
            so.Leave(gm.getGame(cookie).Id)
            gm.joinGame(cookie, host)

            // join new game that the player is in with the other player
            so.Join(gm.getGame(cookie).Id)

            so.Emit("state", gm.getState(cookie))
        })

        so.On("action", func(action string) {
            // TODO handle the action
            actOther := gm.doAction(cookie, action)

            // need check if waiting
            if actOther == "help" {

            } else if actOther == "force" {

            }
            // need check if force
            so.Emit("state", gm.getState(cookie))
        })

        so.Join("chatTest")
        so.On("chat message", func(msg string) {
            log.Println(so.Emit("chat message", msg), " emit:", msg)
            if gm.isInActiveGame(cookie) {
                so.BroadcastTo(gm.getGame(cookie).Id, "chat message", msg)
            } else {
                so.BroadcastTo("chatTest", "chat message", msg)
            }
        })
        so.On("disconnection", func() {
            log.Println(cookie, "disconnected")
        })
    }
}

func socketErrorHandler(so socketio.Socket, err error) {
    log.Println("error:", err)
}
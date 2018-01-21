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
    name (str): sends name so it can be seen client side
    chat message (str): puts a message
    chat util (boolean): true if should display the chat
    room (object): data about who is hosting a game
    state (object): data about the current state of the game
    help (object): data asking if the other player wants to help 
*/

func handleStartup(gm *GameMaster, cookie string, so socketio.Socket) {
    // TODO handle logic for when the game starts and routing rooms and stuffs

    // check if the user has set a name yet
    if gm.getName(cookie) == "" {
        so.Emit("ask name", "I don't know your name yet! :O")
    } else {
        so.Emit("name", gm.getName(cookie))
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
        so.Emit("room", gm.getRooms(cookie))
        time.Sleep(1000 * time.Millisecond)
    }

    log.Println("got here!")
    log.Println(cookie, gm.getState(cookie))
    // emit state information once in a game
    for gm.isInActiveGame(cookie) {
        so.Emit("state", gm.getState(cookie))
        time.Sleep(1000 * time.Millisecond)
    }   
}

// appends the username to the message
func chatMiddleware(gm *GameMaster, cookie string, msg string) string {
    name := gm.getName(cookie)
    if name != "" {
        msg = name + ": " + msg
    }
    return msg
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
            log.Println("Saving player name:", playerName)
            gm.setName(cookie, playerName)
        })

        so.On("join", func(host string) {
            // leave the game that player host alone
            so.Leave(gm.getGame(cookie).Id)
            gm.joinGame(cookie, host)

            // join new game that the player is in with the other player
            so.Join(gm.getGame(cookie).Id)

            // so.Emit("state", gm.getState(cookie))
        })

        so.On("action", func(action string) {
            // actOther := gm.doAction(cookie, action)
            gm.doAction(cookie, action)

            if(!gm.hasEnded(cookie)) {
                // need check if waiting
                // if actOther == "help" {
                //     so.BroadcastTo(gm.getGame(cookie).Id, "help", gm.getState(cookie))
                //     // TODO put timeout and treat it as no
                // } else if actOther == "force" {
                //     so.BroadcastTo(gm.getGame(cookie).Id, "state", gm.getState(cookie))
                // }
                // need check if force
                so.Emit("state", gm.getState(cookie))
            } else {
                so.Emit("end", gm.endState(cookie))
            }

            
        })

        so.Join("chatTest")
        so.On("chat message", func(msg string) {
            msg = chatMiddleware(gm, cookie, msg)
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
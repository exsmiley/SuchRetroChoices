package main

import (
    "log"
    // "fmt"
    "time"
)

type Player struct {
    cookie string
    score int
    states []string
    next string // used after a wait occurs
    ended bool
}

type Game struct {
    Id string // it will be a uuid
    host string
    players map[string]Player // maps player name to map
    waiting string // cookie of player that is waiting for the other one
}

type GameMaster struct {
    playerToGame map[string]string // TODO periodically garbage collect disconnected players to free space
    cookieToName map[string]string
    games map[string]Game
    story *Story
}

type GameRoom struct {
    HostId string
    HostName string
}


type GameRooms struct {
    Rooms []GameRoom
}

// struct to send back and forth
type GameState struct {
    Waiting bool
    Text string
    Special string
    Actions []string
    Image string
}

type EndState struct {
    Text string
    Image string
    Score int
}


func newGameMaster(story *Story) GameMaster {
    gm := GameMaster{}
    gm.playerToGame = make(map[string]string) // cookie name to string
    gm.games = make(map[string]Game)
    gm.cookieToName = make(map[string]string)
    gm.story = story
    return gm
}

// admin game stuff
func (gm *GameMaster) isInActiveGame(playerCookie string) bool {
    game := gm.getGame(playerCookie)
    return len(game.players) == 2
}

func (gm *GameMaster) getGame(playerCookie string) Game {
    return gm.games[gm.playerToGame[playerCookie]]
}

// TODO maybe map cookie to something else for hosts? too lazy for now
// TODO stop race conditions
func (gm *GameMaster) joinGame(playerCookie string, host string) bool {
    game, ok := gm.games[gm.playerToGame[host]]

    if !ok || len(game.players) != 1 {
        return false
    }

    // remove unhosted game
    delete(gm.games, gm.playerToGame[playerCookie])

    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}, "", false}
    gm.playerToGame[playerCookie] = game.Id

    return true
}

func (gm *GameMaster) hostGame(playerCookie string) string {
    // gm.games[gm.playerToGame]
    game := Game{genId(), playerCookie, map[string]Player{}, ""}
    game.players[playerCookie] = Player{playerCookie, 0, []string{"start"}, "", false}

    gm.playerToGame[playerCookie] = game.Id
    gm.games[game.Id] = game
    return game.Id
}

func (gm *GameMaster) getOtherPlayer(playerCookie string) Player {
    game := gm.games[gm.playerToGame[playerCookie]]
    for cookie, player := range game.players {
        if cookie != playerCookie {
            return player
        }
    }
    // should never get here
    return Player{}
}

func (gm *GameMaster) getRooms(playerCookie string) GameRooms {
    grs := GameRooms{[]GameRoom{}}

    hosts := map[string]bool{}

    for _, game := range gm.games {
        hostId := game.host
        if hostId == playerCookie || hosts[hostId] {
            continue
        }

        hostName := gm.getName(hostId)
        if hostName == "" {
            continue
        }

        room := GameRoom{hostId, hostName}
        grs.Rooms = append(grs.Rooms, room)
        hosts[hostId] = true
    }

    return grs
}


// game utils
func (gm *GameMaster) getPlayer(playerCookie string) Player {
    return gm.games[gm.playerToGame[playerCookie]].players[playerCookie]
}

func (gm *GameMaster) getName(playerCookie string) string {
    return gm.cookieToName[playerCookie]
}

func (gm *GameMaster) setName(playerCookie string, name string) {
    gm.cookieToName[playerCookie] = name
}


// game mechanics
// check if the player's game has ended
func (gm *GameMaster) hasEnded(playerCookie string) bool {
    player := gm.getPlayer(playerCookie)
    states := player.states
    lastState := states[len(states)-1]
    return gm.story.hasEnded(lastState)
}

func (gm *GameMaster) endState(playerCookie string) EndState {
    story := gm.story
    player := gm.getPlayer(playerCookie)
    states := player.states
    lastState := states[len(states)-1]


    image := story.getImage(lastState)
    text := story.getText(lastState)

    es := EndState{text, image, player.score}

    return es
}


func (gm *GameMaster) abortWait(playerCookie string, state string) {
    // give them 15 seconds to catch up TODO maybe make longer
    time.Sleep(15000 * time.Millisecond)

    player := gm.getPlayer(playerCookie)
    states := player.states
    lastState := states[len(states)-1]
    game := gm.games[gm.playerToGame[playerCookie]]

    log.Println("trying to abort", state, lastState, game.waiting)

    if state == lastState && game.waiting == playerCookie {
        next, points := gm.story.abortCondition(lastState)
        player.score += points
        player.states = append(states, next)
        game.waiting = ""
        gm.games[gm.playerToGame[playerCookie]] = game
        gm.games[gm.playerToGame[playerCookie]].players[playerCookie] = player
    }
}


func (gm *GameMaster) doAction(playerCookie string, action string) string {
    player := gm.getPlayer(playerCookie)
    states := player.states
    game := gm.games[gm.playerToGame[playerCookie]]
    lastState := states[len(states)-1]

    log.Println(player, action, lastState)

    next, points := gm.story.makeChoice(lastState, action)
    player.score += points

    log.Println("something", next, points)

    // handle if other player is waiting
    if game.waiting != "" {
        otherCookie := game.waiting
        otherPlayer := gm.getPlayer(otherCookie)
        otherLastState := otherPlayer.states[len(otherPlayer.states)-1]
        otherNext := gm.story.checkConditions(otherLastState, otherLastState, next)
        otherPlayer.next = otherNext

        log.Println("other next", otherPlayer)
        if otherNext != "" {
            gm.games[gm.playerToGame[playerCookie]].players[otherCookie] = otherPlayer

            nextNext := gm.story.checkConditions(next, next, otherLastState)
            player.next = nextNext
            log.Println("other next2", nextNext)

            game.waiting = ""
        }

    } else if gm.story.needsToWait(next) {
        log.Println("trying to wait")
        // TODO check other player first to see if they made an action?
        otherPlayer := gm.getOtherPlayer(playerCookie)
        otherLastState := otherPlayer.states[len(otherPlayer.states)-1]
        log.Println(otherLastState, lastState)

        lastLastState := ""
        if len(states) > 1 {
            lastLastState = states[len(states)-2]
        }

        if otherLastState == lastState || otherLastState == lastLastState {
            game.waiting = playerCookie
            go gm.abortWait(playerCookie, next)
        }
    } else {
        player.next = ""
    }

    actionString := ""
    if gm.story.triggersHelp(lastState, action) {
        actionString = "help"
        game.waiting = playerCookie
        gm.games[gm.playerToGame[playerCookie]] = game
    } else if gm.story.triggersForce(lastState, action) {
        actionString = "force"
    }

    player.states = append(states, next)

    gm.games[gm.playerToGame[playerCookie]] = game
    gm.games[gm.playerToGame[playerCookie]].players[playerCookie] = player

    // log.Println("are they the same?", player, gm.getPlayer(playerCookie))


    return actionString
}


func (gm *GameMaster) getState(playerCookie string) GameState {
    story := gm.story
    player := gm.getPlayer(playerCookie)
    states := player.states
    game := gm.games[gm.playerToGame[playerCookie]]
    lastState := states[len(states)-1]
    waiting := game.waiting == playerCookie

    // log.Println("Getting data from", player)


    actions := story.getActions(lastState, player.next)
    image := story.getImage(lastState)
    text := story.getText(lastState)

    // fmt.Println(game.waiting, playerCookie)

    gs := GameState{waiting, text, "", actions, image}

    // fmt.Println("sending", gs)

    return gs
}

